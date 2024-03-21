// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build js

package js

import (
	"errors"
	"sync"
	"syscall/js"

	"github.com/lens-vm/lens/host-go/engine/module"
)

type wRuntime struct {
	webAssembly js.Value
}

var _ module.Runtime = (*wRuntime)(nil)

func New() module.Runtime {
	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface
	webAssembly := js.Global().Get("WebAssembly")
	return &wRuntime{webAssembly}
}

func (rt *wRuntime) NewModule(wasmBytes []byte) (module.Module, error) {
	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/compile_static
	promise := rt.webAssembly.Call("compile", wasmBytes)
	results, err := await(promise)
	if err != nil {
		return nil, err
	}
	return &wModule{
		module:  results[0],
		runtime: rt,
	}, nil
}

type wModule struct {
	module  js.Value
	runtime *wRuntime
}

var _ module.Module = (*wModule)(nil)

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Instance, error) {
	var nextFunction = func() module.MemSize { return 0 }
	importObject := map[string]any{
		"lens": map[string]any{
			"next": js.FuncOf(func(this js.Value, args []js.Value) any {
				return nextFunction()
			}),
		},
	}

	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/instantiate_static
	promise := m.runtime.webAssembly.Call("instantiate", m.module, importObject)
	results, err := await(promise)
	if err != nil {
		return module.Instance{}, err
	}
	instance := results[0].Get("instance")
	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/Instance/exports
	exports := instance.Get("exports")

	return module.Instance{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			result := exports.Call("alloc", u)
			return module.MemSize(result.Int()), nil
		},
		Transform: func(next func() module.MemSize) (module.MemSize, error) {
			// By assigning the next function immediately prior to calling transform, we allow multiple
			// pipeline stages to share the same wasm instance - provided they are not called concurrently.
			// This also allows module state to be shared across pipeline stages.
			nextFunction = next
			result := exports.Call("transform")
			return module.MemSize(result.Int()), nil
		},
		GetData: func() []byte {
			memory := exports.Get("memory")
			buffer := memory.Get("buffer")

			data := make([]byte, buffer.Get("length").Int())
			js.CopyBytesToGo(data, buffer)

			return data
		},
		OwnedBy: instance,
	}, nil
}

// await is a helper that waits for and returns results from the given promise.
func await(promise js.Value) ([]js.Value, error) {
	var (
		res []js.Value
		err error
		wg  sync.WaitGroup
	)

	success := js.FuncOf(func(this js.Value, args []js.Value) any {
		defer wg.Done()
		res = args
		return nil
	})
	defer success.Release()

	failure := js.FuncOf(func(this js.Value, args []js.Value) any {
		defer wg.Done()
		err = errors.New(args[0].Call("toString").String())
		return nil
	})
	defer failure.Release()

	wg.Add(1)
	promise.Call("then", success, failure)
	wg.Wait()

	return res, err
}

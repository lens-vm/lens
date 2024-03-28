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
	// copy bytes to JavaScript value
	wasmBytesJS := js.Global().Get("Uint8Array").New(len(wasmBytes))
	js.CopyBytesToJS(wasmBytesJS, wasmBytes)

	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/compile_static
	promise := rt.webAssembly.Call("compile", wasmBytesJS)

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
	// Register the `lens.next` function required as an import for wasm lens modules
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
	instance := results[0]

	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/Instance/exports
	// exports := instance.Get("exports")

	// https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/Memory
	// memory := exports.Get("memory")
	// if memory.Type() != js.TypeObject {
	// 	return module.Instance{}, errors.New("Export 'memory' does not exist")
	// }

	// alloc := exports.Get("alloc")
	// if alloc.Type() != js.TypeFunction {
	// 	return module.Instance{}, errors.New("Export 'alloc' does not exist")
	// }

	// transform := exports.Get(functionName)
	// if transform.Type() != js.TypeFunction {
	// 	return module.Instance{}, errors.New(fmt.Sprintf("Export '%s' does not exist", functionName))
	// }

	params := map[string]any{}
	// Merge the param sets into a single map in case more than
	// one map is provided.
	for _, paramSet := range paramSets {
		for key, value := range paramSet {
			params[key] = value
		}
	}

	if len(params) > 0 {
		// 	setParam := exports.Get("set_param")
		// 	if setParam.Type() != js.TypeFunction {
		// 		return module.Instance{}, errors.New("Export 'set_param' does not exist")
		// 	}

		// 	sourceBytes, err := json.Marshal(params)
		// 	if err != nil {
		// 		return module.Instance{}, err
		// 	}

		// 	// allocate memory to write to
		// 	index := alloc.Invoke(module.TypeIdSize + module.MemSize(len(sourceBytes)) + module.LenSize)
		// 	// read the JavaScript memory into a go slice
		// 	temp := getData()

		// 	err = pipes.WriteItem(module.JSONTypeID, sourceBytes, temp[index.Int():])
		// 	if err != nil {
		// 		return module.Instance{}, err
		// 	}

		// 	// copy the bytes back to JavaScript memory
		// 	js.CopyBytesToJS(data, temp)
		// 	// set param from JavaScript memory
		// 	r := setParam.Invoke(index)
		// 	// read the JavaScript memory into a go slice
		// 	temp = getData()

		// 	// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// 	// from memory using `pipes.GetItem`.
		// 	_, err = pipes.GetItem(temp, module.MemSize(r.Int()))
		// 	if err != nil {
		// 		return module.Instance{}, err
		// 	}
	}

	return module.Instance{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			result := instance.Get("exports").Call("alloc", int32(u))
			return module.MemSize(result.Int()), nil
		},
		Transform: func(next func() module.MemSize) (module.MemSize, error) {
			// By assigning the next function immediately prior to calling transform, we allow multiple
			// pipeline stages to share the same wasm instance - provided they are not called concurrently.
			// This also allows module state to be shared across pipeline stages.
			nextFunction = next
			result := instance.Get("exports").Call(functionName)
			return module.MemSize(result.Int()), nil
		},
		GetData: func() []byte {
			// GetData should return a byte array that supports both read and write!!!!!

			buffer := instance.Get("exports").Get("memory").Get("buffer")
			data := js.Global().Get("Uint8Array").New(buffer)
			temp := make([]byte, data.Get("length").Int())
			js.CopyBytesToGo(temp, data)
			return temp
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

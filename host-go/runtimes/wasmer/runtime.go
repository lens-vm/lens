// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !windows && !js

package wasmer

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type wRuntime struct {
	store *wasmer.Store
}

var _ module.Runtime = (*wRuntime)(nil)

func New() module.Runtime {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	return &wRuntime{
		store: store,
	}
}

type wModule struct {
	runtime *wRuntime
	module  *wasmer.Module
}

var _ module.Module = (*wModule)(nil)

func (rt *wRuntime) NewModule(wasmBytes []byte) (module.Module, error) {
	module, err := wasmer.NewModule(rt.store, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &wModule{
		runtime: rt,
		module:  module,
	}, nil
}

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Instance, error) {
	importObject := wasmer.NewImportObject()

	var nextFunction = func() module.MemSize { return 0 }
	// Register the `lens.next` function required as an import for wasm lens modules
	importObject.Register(
		"lens",
		map[string]wasmer.IntoExtern{
			"next": wasmer.NewFunction(
				m.runtime.store,
				wasmer.NewFunctionType(
					wasmer.NewValueTypes(),
					// Warning: wasmer requires a concrete type here and as such this line is coupled to the module's runtime
					wasmer.NewValueTypes(wasmer.I32),
				),
				func(v []wasmer.Value) ([]wasmer.Value, error) {
					r := nextFunction()
					return []wasmer.Value{wasmer.NewI32(r)}, nil
				},
			),
		},
	)

	instance, err := wasmer.NewInstance(m.module, importObject)
	if err != nil {
		return module.Instance{}, err
	}

	memory, err := instance.Exports.GetMemory("memory")
	if err != nil {
		return module.Instance{}, err
	}

	alloc, err := instance.Exports.GetRawFunction("alloc")
	if err != nil {
		return module.Instance{}, err
	}

	transform, err := instance.Exports.GetRawFunction(functionName)
	if err != nil {
		return module.Instance{}, err
	}

	params := map[string]any{}
	// Merge the param sets into a single map in case more than
	// one map is provided.
	for _, paramSet := range paramSets {
		for key, value := range paramSet {
			params[key] = value
		}
	}

	if len(params) > 0 {
		setParam, err := instance.Exports.GetRawFunction("set_param")
		if err != nil {
			return module.Instance{}, err
		}

		sourceBytes, err := json.Marshal(params)
		if err != nil {
			return module.Instance{}, err
		}

		index, err := alloc.Call(module.TypeIdSize + module.MemSize(len(sourceBytes)) + module.LenSize)
		if err != nil {
			return module.Instance{}, err
		}

		err = pipes.WriteItem(module.JSONTypeID, sourceBytes, memory.Data()[index.(module.MemSize):])
		if err != nil {
			return module.Instance{}, err
		}

		r, err := setParam.Call(index)
		if err != nil {
			return module.Instance{}, err
		}

		// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// from memory using `pipes.GetItem`.
		_, err = pipes.GetItem(memory.Data(), r.(module.MemSize))
		if err != nil {
			return module.Instance{}, err
		}
	}

	return module.Instance{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			r, err := alloc.Call(u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		Transform: func(next func() module.MemSize) (module.MemSize, error) {
			// By assigning the next function immediately prior to calling transform, we allow multiple
			// pipeline stages to share the same wasm instance - provided they are not called concurrently.
			// This also allows module state to be shared across pipeline stages.
			nextFunction = next
			r, err := transform.Call()
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		GetData: memory.Data,
		OwnedBy: instance,
	}, nil
}

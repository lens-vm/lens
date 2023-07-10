// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wasmer

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"
	"github.com/lens-vm/lens/host-go/engine/runtime"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type wRuntime struct {
	store *wasmer.Store
}

var _ runtime.Runtime = (*wRuntime)(nil)

func New() runtime.Runtime {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	return &wRuntime{
		store: store,
	}
}

type wModule struct {
	module *wasmer.Module
}

var _ runtime.Module = (*wModule)(nil)

func (rt *wRuntime) NewModule(wasmBytes []byte) (runtime.Module, error) {
	module, err := wasmer.NewModule(rt.store, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &wModule{
		module: module,
	}, nil
}

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Module, error) {
	importObject := wasmer.NewImportObject()
	instance, err := wasmer.NewInstance(m.module, importObject)
	if err != nil {
		return module.Module{}, err
	}

	memory, err := instance.Exports.GetMemory("memory")
	if err != nil {
		return module.Module{}, err
	}

	alloc, err := instance.Exports.GetRawFunction("alloc")
	if err != nil {
		return module.Module{}, err
	}

	transform, err := instance.Exports.GetRawFunction(functionName)
	if err != nil {
		return module.Module{}, err
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
			return module.Module{}, err
		}

		sourceBytes, err := json.Marshal(params)
		if err != nil {
			return module.Module{}, err
		}

		index, err := alloc.Call(module.TypeIdSize + module.MemSize(len(sourceBytes)) + module.LenSize)
		if err != nil {
			return module.Module{}, err
		}

		err = pipes.WriteItem(module.JSONTypeID, sourceBytes, memory.Data()[index.(module.MemSize):])
		if err != nil {
			return module.Module{}, err
		}

		r, err := setParam.Call(index)
		if err != nil {
			return module.Module{}, err
		}

		// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// from memory using `pipes.GetItem`.
		_, err = pipes.GetItem(memory.Data(), r.(module.MemSize))
		if err != nil {
			return module.Module{}, err
		}
	}

	return module.Module{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			r, err := alloc.Call(u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		Transform: func(u module.MemSize) (module.MemSize, error) {
			r, err := transform.Call(u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		GetData: memory.Data,
		OwnedBy: instance,
	}, nil
}

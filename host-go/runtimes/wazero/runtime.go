// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wazero

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"
	"github.com/tetratelabs/wazero"
)

type wRuntime struct {
	runtime wazero.Runtime
}

var _ module.Runtime = (*wRuntime)(nil)

func New() module.Runtime {
	ctx := context.TODO()
	return &wRuntime{
		runtime: wazero.NewRuntime(ctx),
	}
}

type wModule struct {
	runtime wazero.Runtime
	module  wazero.CompiledModule
}

var _ module.Module = (*wModule)(nil)

func (rt *wRuntime) NewModule(wasmBytes []byte) (module.Module, error) {
	ctx := context.TODO()
	compiledWasm, err := rt.runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &wModule{
		runtime: rt.runtime,
		module:  compiledWasm,
	}, nil
}

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Instance, error) {
	ctx := context.TODO()
	instance, err := m.runtime.InstantiateModule(ctx, m.module, wazero.NewModuleConfig().WithName(""))

	memory := instance.ExportedMemory("memory")
	if memory == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "memory"))
	}

	alloc := instance.ExportedFunction("alloc")
	if alloc == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "alloc"))
	}

	transform := instance.ExportedFunction(functionName)
	if transform == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", functionName))
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
		setParam := instance.ExportedFunction("set_param")
		if err != nil {
			return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "set_param"))
		}

		sourceBytes, err := json.Marshal(params)
		if err != nil {
			return module.Instance{}, err
		}

		index, err := alloc.Call(ctx, uint64(module.TypeIdSize+module.MemSize(len(sourceBytes))+module.LenSize))
		if err != nil {
			return module.Instance{}, err
		}

		data, _ := memory.Read(0, memory.Size())
		err = pipes.WriteItem(module.JSONTypeID, sourceBytes, data[index[0]:])
		if err != nil {
			return module.Instance{}, err
		}

		r, err := setParam.Call(ctx, index[0])
		if err != nil {
			return module.Instance{}, err
		}

		data, _ = memory.Read(0, memory.Size())
		// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// from memory using `pipes.GetItem`.
		_, err = pipes.GetItem(data, module.MemSize(r[0]))
		if err != nil {
			return module.Instance{}, err
		}
	}

	return module.Instance{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			r, err := alloc.Call(ctx, uint64(u))
			if err != nil {
				return 0, err
			}
			return module.MemSize(r[0]), nil
		},
		Transform: func(u module.MemSize) (module.MemSize, error) {
			r, err := transform.Call(ctx, uint64(u))
			if err != nil {
				return 0, err
			}
			return module.MemSize(r[0]), nil
		},
		GetData: func() []byte {
			data, _ := memory.Read(0, memory.Size())
			return data
		},
		OwnedBy: instance,
	}, nil
}

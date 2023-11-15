// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package wasmtime

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"

	"github.com/bytecodealliance/wasmtime-go/v14"
)

type wRuntime struct {
	store *wasmtime.Store
}

var _ module.Runtime = (*wRuntime)(nil)

func New() module.Runtime {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	wasiConfig := wasmtime.NewWasiConfig()
	store.SetWasi(wasiConfig)

	return &wRuntime{
		store: store,
	}
}

type wModule struct {
	rt     *wRuntime
	module *wasmtime.Module
	wasi   bool
}

var _ module.Module = (*wModule)(nil)

func (rt *wRuntime) NewModule(wasmBytes []byte, wasi bool) (module.Module, error) {
	module, err := wasmtime.NewModule(rt.store.Engine, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &wModule{
		rt:     rt,
		module: module,
		wasi:   wasi,
	}, nil
}

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Instance, error) {
	var instance *wasmtime.Instance
	var err error

	switch {
	case m.wasi:
		// wasi imports
		linker := wasmtime.NewLinker(m.rt.store.Engine)
		if err := linker.DefineWasi(); err != nil {
			return module.Instance{}, err
		}

		instance, err = linker.Instantiate(m.rt.store, m.module)
		if err != nil {
			return module.Instance{}, err
		}
	default:
		// default imports
		instance, err = wasmtime.NewInstance(m.rt.store, m.module, []wasmtime.AsExtern{})
		if err != nil {
			return module.Instance{}, err
		}
	}

	mem := instance.GetExport(m.rt.store, "memory")
	if mem == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "memory"))
	}

	memory := mem.Memory()
	if memory == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "memory"))
	}

	alloc := instance.GetFunc(m.rt.store, "alloc")
	if alloc == nil {
		return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "alloc"))
	}

	transform := instance.GetFunc(m.rt.store, functionName)
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
		setParam := instance.GetFunc(m.rt.store, "set_param")
		if setParam == nil {
			return module.Instance{}, errors.New(fmt.Sprintf("Export `%s` does not exist", "set_param"))
		}

		sourceBytes, err := json.Marshal(params)
		if err != nil {
			return module.Instance{}, err
		}

		index, err := alloc.Call(m.rt.store, module.TypeIdSize+module.MemSize(len(sourceBytes))+module.LenSize)
		if err != nil {
			return module.Instance{}, err
		}

		err = pipes.WriteItem(module.JSONTypeID, sourceBytes, memory.UnsafeData(m.rt.store)[index.(module.MemSize):])
		if err != nil {
			return module.Instance{}, err
		}

		r, err := setParam.Call(m.rt.store, index)
		if err != nil {
			return module.Instance{}, err
		}

		// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// from memory using `pipes.GetItem`.
		_, err = pipes.GetItem(memory.UnsafeData(m.rt.store), r.(module.MemSize))
		if err != nil {
			return module.Instance{}, err
		}
	}

	return module.Instance{
		Alloc: func(u module.MemSize) (module.MemSize, error) {
			r, err := alloc.Call(m.rt.store, u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		Transform: func(u module.MemSize) (module.MemSize, error) {
			r, err := transform.Call(m.rt.store, u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		GetData: func() []byte { return memory.UnsafeData(m.rt.store) },
		OwnedBy: instance,
	}, nil
}

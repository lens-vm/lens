// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !js

package wazero

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"
	"github.com/tetratelabs/wazero"
)

type wRuntime struct {
	compilationCache wazero.CompilationCache
}

var _ module.Runtime = (*wRuntime)(nil)

// New creates a new wazero wasm runtime.
//
// WARNING: This runtime does not current allow for instance reuse within a single pipeline.
// Please see https://github.com/lens-vm/lens/issues/71 for more info.
func New() module.Runtime {
	return &wRuntime{
		compilationCache: wazero.NewCompilationCache(),
	}
}

type wModule struct {
	compilationCache wazero.CompilationCache
	moduleBytes      []byte
}

var _ module.Module = (*wModule)(nil)

func (rt *wRuntime) NewModule(wasmBytes []byte) (module.Module, error) {
	return &wModule{
		compilationCache: rt.compilationCache,
		moduleBytes:      wasmBytes,
	}, nil
}

func (m *wModule) NewInstance(functionName string, paramSets ...map[string]any) (module.Instance, error) {
	ctx := context.TODO()
	runtimeConfig := wazero.NewRuntimeConfig().WithCompilationCache(m.compilationCache)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	var nextFunction = func() module.MemSize { return 0 }
	_, err := runtime.NewHostModuleBuilder("lens").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) module.MemSize {
			return nextFunction()
		}).
		Export("next").
		Instantiate(ctx)
	if err != nil {
		return module.Instance{}, err
	}

	instance, err := runtime.Instantiate(ctx, m.moduleBytes)
	if err != nil {
		return module.Instance{}, err
	}

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

		mem := newMemory(memory, int32(index[0]))
		err = pipes.WriteItem(mem, module.JSONTypeID, sourceBytes)
		if err != nil {
			return module.Instance{}, err
		}

		r, err := setParam.Call(ctx, index[0])
		if err != nil {
			return module.Instance{}, err
		}

		// The `set_param` wasm function may error, in which case the error needs to be retrieved
		// from memory using `pipes.GetItem`.
		mem = newMemory(memory, int32(r[0]))
		_, err = pipes.ReadItem(mem)
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
		Transform: func(next func() module.MemSize) (module.MemSize, error) {
			// By assigning the next function immediately prior to calling transform, we allow multiple
			// pipeline stages to share the same wasm instance - provided they are not called concurrently.
			// This also allows module state to be shared across pipeline stages.
			nextFunction = next
			r, err := transform.Call(ctx)
			if err != nil {
				return 0, err
			}
			return module.MemSize(r[0]), nil
		},
		Memory: func(offset int32) io.ReadWriter {
			return newMemory(memory, offset)
		},
		OwnedBy: instance,
	}, nil
}

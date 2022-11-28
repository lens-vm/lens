package engine

import (
	"encoding/json"
	"os"

	"github.com/lens-vm/lens/host-go/engine/internal/pipes"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/sourcenetwork/immutable/enumerable"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// Append appends the given Module(s) to the given source Enumerable, returning the result.
//
// It will try and find the optimal way to communicate between the source and the new module, returning an enumerable of a type
// that best fits the situation. The source can be any type that implements the Enumerable interface, it does not need to be a
// lens module.
func Append[TSource any, TResult any](src enumerable.Enumerable[TSource], modules ...module.Module) enumerable.Enumerable[TResult] {
	if len(modules) == 0 {
		return src.(enumerable.Enumerable[TResult])
	}

	if len(modules) == 1 {
		return append[TSource, TResult](src, modules[0])
	}

	intermediarySource := append[TSource, map[string]any](src, modules[0])
	for i := 1; i < len(modules)-1; i++ {
		intermediarySource = append[map[string]any, map[string]any](intermediarySource, modules[i])
	}

	return append[map[string]any, TResult](intermediarySource, modules[len(modules)-1])
}

func append[TSource any, TResult any](src enumerable.Enumerable[TSource], module module.Module) enumerable.Enumerable[TResult] {
	switch typedSrc := src.(type) {
	case pipes.Pipe[TSource]:
		return pipes.NewFromPipe[TSource, TResult](typedSrc, module)
	default:
		return pipes.NewFromSource[TSource, TResult](src, module)
	}
}

// LoadModule loads a lens at the given path.
func LoadModule(path string, paramSets ...map[string]any) (module.Module, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return module.Module{}, err
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	wasmModule, err := wasmer.NewModule(store, content)
	if err != nil {
		return module.Module{}, err
	}

	importObject := wasmer.NewImportObject()
	instance, err := wasmer.NewInstance(wasmModule, importObject)
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

	transform, err := instance.Exports.GetRawFunction("transform")
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

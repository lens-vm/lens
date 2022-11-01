package lib

import (
	"io/ioutil"
	"lens-host/lib/enumerable"
	"lens-host/lib/internal/pipes"
	"lens-host/lib/module"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// AppendLens appends the given Module to the given source Enumerable, returning the result.
//
// It will try and find the optimal way to communicate between the source and the new module, returning an enumerable of a type
// that best fits the situation. The source can be any type that implements the Enumerable interface, it does not need to be a
// lens module.
func AppendLens[TSource any, TResult any](src enumerable.Enumerable[TSource], module module.Module) enumerable.Enumerable[TResult] {
	switch typedSrc := src.(type) {
	case pipes.Pipe[TSource]:
		return pipes.FromModuleAsFullToFull[TSource, TResult](typedSrc, module)
	default:
		return pipes.FromSourceToFull[TSource, TResult](src, module)
	}
}

// LoadModule loads a lens at the given path.
func LoadModule(path string) (module.Module, error) {
	content, err := ioutil.ReadFile(path)
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

	return module.New(
		func(u module.MemSize) (module.MemSize, error) {
			r, err := alloc.Call(u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		func(u module.MemSize, a ...any) (module.MemSize, error) {
			r, err := transform.Call(u)
			if err != nil {
				return 0, err
			}
			return r.(module.MemSize), err
		},
		memory.Data,
		instance,
	), nil
}

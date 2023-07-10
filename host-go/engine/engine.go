// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package engine

import (
	"os"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/engine/pipes"
	"github.com/lens-vm/lens/host-go/engine/runtime"
	"github.com/sourcenetwork/immutable/enumerable"
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

func NewModule(runtime runtime.Runtime, path string) (runtime.Module, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return runtime.NewModule(content)
}

func NewInstance(module runtime.Module, paramSets ...map[string]any) (module.Module, error) {
	return module.NewInstance("transform", paramSets...)
}

func NewInverse(module runtime.Module, paramSets ...map[string]any) (module.Module, error) {
	return module.NewInstance("inverse", paramSets...)
}

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

// Append appends the given Module Instances to the given source Enumerable, returning the result.
//
// It will try and find the optimal way to communicate between the source and the new module instance, returning an enumerable of a type
// that best fits the situation. The source can be any type that implements the Enumerable interface, it does not need to be a
// lens module instance.
func Append[TSource any, TResult any](src enumerable.Enumerable[TSource], instances ...module.Instance) enumerable.Enumerable[TResult] {
	if len(instances) == 0 {
		return src.(enumerable.Enumerable[TResult])
	}

	if len(instances) == 1 {
		return append[TSource, TResult](src, instances[0])
	}

	intermediarySource := append[TSource, map[string]any](src, instances[0])
	for i := 1; i < len(instances)-1; i++ {
		intermediarySource = append[map[string]any, map[string]any](intermediarySource, instances[i])
	}

	return append[map[string]any, TResult](intermediarySource, instances[len(instances)-1])
}

func append[TSource any, TResult any](src enumerable.Enumerable[TSource], instance module.Instance) enumerable.Enumerable[TResult] {
	switch typedSrc := src.(type) {
	case pipes.Pipe[TSource]:
		return pipes.NewFromPipe[TSource, TResult](typedSrc, instance)
	default:
		return pipes.NewFromSource[TSource, TResult](src, instance)
	}
}

func NewModule(runtime runtime.Runtime, path string) (runtime.Module, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return runtime.NewModule(content)
}

func NewInstance(module runtime.Module, paramSets ...map[string]any) (module.Instance, error) {
	return module.NewInstance("transform", paramSets...)
}

func NewInverse(module runtime.Module, paramSets ...map[string]any) (module.Instance, error) {
	return module.NewInstance("inverse", paramSets...)
}

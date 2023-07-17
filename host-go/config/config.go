// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/lens-vm/lens/host-go/config/internal/json"
	"github.com/lens-vm/lens/host-go/config/model"
	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wasmer"
	"github.com/sourcenetwork/immutable/enumerable"
)

// LoadFromFile loads a lens file at the given path and applies it to the provided src.
//
// It does not enumerate the src.
func LoadFromFile[TSource any, TResult any](path string, src enumerable.Enumerable[TSource]) (enumerable.Enumerable[TResult], error) {
	// We only support json lens files at the moment, so we just trust that it is json.
	// In the future we'll need to determine which format the file is in.
	lensConfig, err := json.Load(path)
	if err != nil {
		return nil, err
	}

	return Load[TSource, TResult](lensConfig, src)
}

// Load constructs a lens from the given config and applies it to the provided src.
//
// It does not enumerate the src.
func Load[TSource any, TResult any](lensConfig model.Lens, src enumerable.Enumerable[TSource]) (enumerable.Enumerable[TResult], error) {
	runtime := wasmer.New()

	modulesByPath := map[string]module.Module{}
	for _, moduleCfg := range lensConfig.Lenses {
		if _, ok := modulesByPath[moduleCfg.Path]; ok {
			continue
		}

		// Modules are fairly expensive objects, and they can be reused, so we de-duplicate
		// the WAT code paths here and make sure we only create unique module objects.
		lensModule, err := engine.NewModule(runtime, moduleCfg.Path)
		if err != nil {
			return nil, err
		}
		modulesByPath[moduleCfg.Path] = lensModule
	}

	instances := []module.Instance{}
	for _, moduleCfg := range lensConfig.Lenses {
		lensModule := modulesByPath[moduleCfg.Path]

		var instance module.Instance
		var err error
		if moduleCfg.Inverse {
			instance, err = engine.NewInverse(lensModule, moduleCfg.Arguments)
		} else {
			instance, err = engine.NewInstance(lensModule, moduleCfg.Arguments)
		}

		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return engine.Append[TSource, TResult](src, instances...), nil
}

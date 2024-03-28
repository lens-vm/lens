// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/lens-vm/lens/host-go/config/model"
	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes"
	"github.com/sourcenetwork/immutable/enumerable"
)

// Load constructs a lens from the given config and applies it to the provided src.
//
// It does not enumerate the src.
func Load[TSource any, TResult any](lensConfig model.Lens, src enumerable.Enumerable[TSource]) (enumerable.Enumerable[TResult], error) {
	runtime := runtimes.Default()
	modulesByHash := map[string]module.Module{}

	return LoadInto[TSource, TResult](runtime, modulesByHash, lensConfig, src)
}

// LoadInto constructs a lens from the given config and applies it to the provided src
// extending the provided runtime and module cache.
//
// It does not enumerate the src. Any new modules will be added to the given module map.
func LoadInto[TSource any, TResult any](
	runtime module.Runtime,
	modulesByHash map[string]module.Module,
	lensConfig model.Lens,
	src enumerable.Enumerable[TSource],
) (enumerable.Enumerable[TResult], error) {
	for _, moduleCfg := range lensConfig.Lenses {
		hash := moduleCfg.Hash()
		// Modules are fairly expensive objects, and they can be reused, so we de-duplicate
		// the WAT code paths here and make sure we only create unique module objects.
		if _, ok := modulesByHash[hash]; ok {
			continue
		}

		lensModule, err := engine.NewModule(runtime, moduleCfg.Content)
		if err != nil {
			return nil, err
		}
		modulesByHash[hash] = lensModule
	}

	instances := []module.Instance{}
	for _, moduleCfg := range lensConfig.Lenses {
		hash := moduleCfg.Hash()
		lensModule := modulesByHash[hash]

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

package config

import (
	"github.com/lens-vm/lens/host-go/config/internal/json"
	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/enumerable"
	"github.com/lens-vm/lens/host-go/engine/module"
)

// Load loads a lens file at the given path and applies it to the provided src.
//
// It does not enumerate the src.
func Load[TSource any, TResult any](path string, src enumerable.Enumerable[TSource]) (enumerable.Enumerable[TResult], error) {
	// We only support json lens files at the moment, so we just trust that it is json.
	// In the future we'll need to determine which format the file is in.
	lensConfig, err := json.Load(path)
	if err != nil {
		return nil, err
	}

	modules := []module.Module{}
	for _, lensModule := range lensConfig.Lenses {
		module, err := engine.LoadModule(lensModule.Path, lensModule.AdditionalParameters...)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	return engine.Append[TSource, TResult](src, modules...), nil
}

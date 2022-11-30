package json

import (
	"encoding/json"
	"os"

	"github.com/lens-vm/lens/host-go/config/internal/model"
)

type Lens struct {
	Lenses []LensModule `json:"lenses"`
}

type LensModule struct {
	Path      string         `json:"path"`
	Inverse   bool           `json:"inverse"`
	Arguments map[string]any `json:"arguments"`
}

func Load(path string) (model.Lens, error) {
	lensFileJson, err := os.ReadFile(path)
	if err != nil {
		return model.Lens{}, err
	}

	var lensFile Lens
	err = json.Unmarshal(lensFileJson, &lensFile)
	if err != nil {
		return model.Lens{}, err
	}

	lenses := make([]model.LensModule, len(lensFile.Lenses))
	for i, lensModule := range lensFile.Lenses {
		lenses[i] = model.LensModule{
			Path:      lensModule.Path,
			Inverse:   lensModule.Inverse,
			Arguments: lensModule.Arguments,
		}
	}

	return model.Lens{
		Lenses: lenses,
	}, nil
}

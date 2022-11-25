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
	Path      string `json:"path"`
	Arguments []any  `json:"arguments"`
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
		arguments := make([]any, len(lensModule.Arguments))
		for j, additionalParameter := range lensModule.Arguments {
			arguments[j] = additionalParameter
		}

		lenses[i] = model.LensModule{
			Path:      lensModule.Path,
			Arguments: arguments,
		}
	}

	return model.Lens{
		Lenses: lenses,
	}, nil
}

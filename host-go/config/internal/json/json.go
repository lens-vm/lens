package json

import (
	"encoding/json"
	"io/ioutil"

	"github.com/lens-vm/lens/host-go/config/internal/model"
)

type Lens struct {
	Lenses []LensModule `json:"lenses"`
}

type LensModule struct {
	Path                 string `json:"path"`
	AdditionalParameters []any  `json:"additionalParameters"`
}

func Load(path string) (model.Lens, error) {
	lensFileJson, err := ioutil.ReadFile(path)
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
		additionalParameters := make([]any, len(lensModule.AdditionalParameters))
		for j, additionalParameter := range lensModule.AdditionalParameters {
			additionalParameters[j] = additionalParameter
		}

		lenses[i] = model.LensModule{
			Path:                 lensModule.Path,
			AdditionalParameters: additionalParameters,
		}
	}

	return model.Lens{
		Lenses: lenses,
	}, nil
}

package integration

import (
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestWithParams(t *testing.T) {
	type Input struct {
		Name string
		Age  int
	}

	type Output struct {
		MiddleName string
		Age        int
	}

	executeTest(
		t,
		TestCase[Input, Output]{
			LensFile: `
			{
				"lenses": [
					{
						"path": "` + modules.WasmPath4 + `",
						"additionalParameters": [
							"Name",
							"MiddleName"
						]
					}
				]
			}`,
			Input: []Input{
				{
					Name: "John",
					Age:  3,
				},
			},
			ExpectedOutput: []Output{
				{
					MiddleName: "John",
					Age:        3,
				},
			},
		},
	)
}

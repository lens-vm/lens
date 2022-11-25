package integration

import (
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestSimple(t *testing.T) {
	type Input struct {
		Name string
		Age  int
	}

	type Output struct {
		FullName string
		Age      int
	}

	executeTest(
		t,
		TestCase[Input, Output]{
			LensFile: `
			{
				"lenses": [
					{
						"path": "` + modules.WasmPath1 + `"
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
					FullName: "John",
					Age:      3,
				},
			},
		},
	)
}

package integration

import (
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestSimpleWithModules(t *testing.T) {
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
					},
					{
						"path": "` + modules.WasmPath2 + `"
					}
				]
			}`,
			Input: []Input{
				{
					Name: "John",
					Age:  3,
				},
				{
					Name: "Fred",
					Age:  5,
				},
				{
					Name: "Orpheus",
					Age:  7,
				},
			},
			ExpectedOutput: []Output{
				{
					FullName: "John",
					Age:      4,
				},
				{
					FullName: "Fred",
					Age:      6,
				},
				{
					FullName: "Orpheus",
					Age:      8,
				},
			},
		},
	)
}

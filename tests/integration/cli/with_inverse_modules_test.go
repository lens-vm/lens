package integration

import (
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestInverseWithModules(t *testing.T) {
	type Value struct {
		FullName string
		Age      int
	}

	executeTest(
		t,
		TestCase[Value, Value]{
			LensFile: `
			{
				"lenses": [
					{
						"path": "` + modules.WasmPath2 + `"
					},
					{
						"path": "` + modules.WasmPath2 + `",
						"inverse": true
					},
					{
						"path": "` + modules.WasmPath2 + `",
						"inverse": true
					}
				]
			}`,
			Input: []Value{
				{
					FullName: "John",
					Age:      3,
				},
				{
					FullName: "Fred",
					Age:      5,
				},
				{
					FullName: "Orpheus",
					Age:      7,
				},
			},
			ExpectedOutput: []Value{
				{
					FullName: "John",
					Age:      2,
				},
				{
					FullName: "Fred",
					Age:      4,
				},
				{
					FullName: "Orpheus",
					Age:      6,
				},
			},
		},
	)
}

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
		},
	)
}

func TestSimpleWithEmptyItem(t *testing.T) {
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
				{},
				{
					Name: "Orpheus",
					Age:  7,
				},
			},
			ExpectedOutput: []Output{
				{
					FullName: "John",
					Age:      3,
				},
				{},
				{
					FullName: "Orpheus",
					Age:      7,
				},
			},
		},
	)
}

func TestSimpleWithNilItem(t *testing.T) {
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
		TestCase[*Input, *Output]{
			LensFile: `
			{
				"lenses": [
					{
						"path": "` + modules.WasmPath1 + `"
					}
				]
			}`,
			Input: []*Input{
				{
					Name: "John",
					Age:  3,
				},
				nil,
				{
					Name: "Orpheus",
					Age:  7,
				},
			},
			ExpectedOutput: []*Output{
				{
					FullName: "John",
					Age:      3,
				},
				nil,
				{
					FullName: "Orpheus",
					Age:      7,
				},
			},
		},
	)
}

package integration

import (
	"encoding/base64"
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestWithFilter(t *testing.T) {
	type Value struct {
		Name string
		Type string `json:"__type"`
	}

	executeTest(
		t,
		TestCase[Value, Value]{
			LensFile: `
			{
				"lenses": [
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Filter) + `"
					}
				]
			}`,
			Input: []Value{
				{
					Name: "John",
					Type: "pass",
				},
				{
					Name: "Fred",
					Type: "skip",
				},
				{
					Name: "Orpheus",
					Type: "pass",
				},
			},
			ExpectedOutput: []Value{
				{
					Name: "John",
					Type: "pass",
				},
				{
					Name: "Orpheus",
					Type: "pass",
				},
			},
		},
	)
}

func TestWithNormalize(t *testing.T) {
	type Book struct {
		Name        string
		PageNumbers []int32
	}
	type Page struct {
		BookName string
		Number   int32
	}

	executeTest(
		t,
		TestCase[Book, Page]{
			LensFile: `
			{
				"lenses": [
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Normalize) + `"
					}
				]
			}`,
			Input: []Book{
				{
					Name:        "The Tiger who came to tea",
					PageNumbers: []int32{1, 2},
				},
				{
					Name:        "The Elephant and the Balloon",
					PageNumbers: []int32{157, 235, 384},
				},
			},
			ExpectedOutput: []Page{
				{
					BookName: "The Tiger who came to tea",
					Number:   1,
				},
				{
					BookName: "The Tiger who came to tea",
					Number:   2,
				},
				{
					BookName: "The Elephant and the Balloon",
					Number:   157,
				},
				{
					BookName: "The Elephant and the Balloon",
					Number:   235,
				},
				{
					BookName: "The Elephant and the Balloon",
					Number:   384,
				},
			},
		},
	)
}

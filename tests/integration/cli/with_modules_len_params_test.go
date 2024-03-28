package integration

import (
	"encoding/base64"
	"testing"

	"github.com/lens-vm/lens/tests/modules"
)

func TestWithModulesWithNormalizeWithParams(t *testing.T) {
	type Book struct {
		Name        string
		PageNumbers []int32
	}
	type Page struct {
		BookName   string
		PageNumber int32
	}

	executeTest(
		t,
		TestCase[Book, Page]{
			LensFile: `
			{
				"lenses": [
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Normalize) + `"
					},
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Rename) + `",
						"arguments": {
							"src": "Number",
							"dst": "PageNumber"
						}
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
					BookName:   "The Tiger who came to tea",
					PageNumber: 1,
				},
				{
					BookName:   "The Tiger who came to tea",
					PageNumber: 2,
				},
				{
					BookName:   "The Elephant and the Balloon",
					PageNumber: 157,
				},
				{
					BookName:   "The Elephant and the Balloon",
					PageNumber: 235,
				},
				{
					BookName:   "The Elephant and the Balloon",
					PageNumber: 384,
				},
			},
		},
	)
}

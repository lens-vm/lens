// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package integration

import (
	"encoding/base64"
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
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Simple2) + `"
					},
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Simple2) + `",
						"inverse": true
					},
					{
						"content": "` + base64.StdEncoding.EncodeToString(modules.RustWasm32Simple2) + `",
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

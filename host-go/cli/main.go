// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"os"

	"github.com/lens-vm/lens/host-go/config"
	"github.com/lens-vm/lens/host-go/config/model"

	"github.com/sourcenetwork/immutable/enumerable"
)

func main() {
	dec := json.NewDecoder(os.Stdin)

	var lensConfig model.Lens
	err := dec.Decode(&lensConfig)
	if err != nil {
		panic(err)
	}

	var data []map[string]any
	err = dec.Decode(&data)
	if err != nil {
		panic(err)
	}

	src := enumerable.New(data)
	result, err := config.Load[map[string]any, map[string]any](lensConfig, src)
	if err != nil {
		panic(err)
	}

	results := []any{}
	for {
		hasNext, err := result.Next()
		if err != nil {
			panic(err)
		}

		if !hasNext {
			break
		}

		val, err := result.Value()
		if err != nil {
			panic(err)
		}

		results = append(results, val)
	}

	resultJson, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	os.Stdout.WriteString(string(resultJson))
}

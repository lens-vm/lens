// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/sourcenetwork/immutable/enumerable"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppendLensWithoutWasm asserts that AppendLens can function independently of anything wasm related.
// This may be very important at a later date, as it means this function can (in-theory, minus any internal
// dependencies) can be executed in a wasm/wasi runtime. It should also help flag any changes to AppendLens'
// externally visible behaviour that may otherwise be hidden (e.g. if otherwise testing using wasm modules).
//
// TestAppendLensWithoutWasm also helps document how AppendLens works.
//
// It is not recomended to actually use AppendLens like this - if you are not transforming data via a wasm module,
// you should probably be using enumerable.Select or similar instead.
func TestAppendLensWithoutWasm(t *testing.T) {
	sourceSlice := []type1{
		{
			Name: "John",
			Age:  32,
		},
		{
			Name: "Fred",
			Age:  55,
		},
	}
	source := enumerable.New(sourceSlice)

	testModule := newNativeModule()
	results := engine.Append[type1, type2](
		source,
		module.Instance{
			Alloc: func(size module.MemSize) (module.MemSize, error) {
				var arbitraryIndex module.MemSize = 5
				return arbitraryIndex, nil
			},
			Transform: func(startIndex module.MemSize) (module.MemSize, error) {
				return testModule.nativeTransform(startIndex)
			},
			GetData: func() []byte {
				return testModule.memory
			},
		},
	)

	hasNext, err := results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := results.Value()
	require.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      32,
	}, val)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err = results.Value()
	require.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "Fred",
		Age:      55,
	}, val)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

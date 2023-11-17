// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/tests/modules"
	"github.com/sourcenetwork/immutable/enumerable"
	"github.com/stretchr/testify/assert"
)

func TestAllocAndFreeAreCalledCorrectlyPerItemFromSource(t *testing.T) {
	sourceSlice := []type1{
		{
			Name: "John",
			Age:  32,
		},
		{
			Name: "Fred",
			Age:  55,
		},
		{
			Name: "Islam",
			Age:  35,
		},
	}
	source := enumerable.New(sourceSlice)

	timesAllocCalled := 0
	timesFreeCalled := 0
	var lastItemIndex module.MemSize
	var lastItemSize module.MemSize

	testModule := newNativeModule()
	results := engine.Append[type1, type2](
		source,
		module.Instance{
			Alloc: func(size module.MemSize) (module.MemSize, error) {
				lastItemIndex++
				timesAllocCalled++
				lastItemSize = size
				return lastItemIndex, nil
			},
			Free: func(ptr, size module.MemSize) error {
				timesFreeCalled++
				// Assert that we are freeing the last allocated item
				assert.Equal(t, ptr, lastItemIndex)
				assert.Equal(t, size, lastItemSize)
				return nil
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

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)

	// Assert alloc and free were called once per item
	assert.Equal(t, len(sourceSlice), timesAllocCalled)
	assert.Equal(t, len(sourceSlice), timesFreeCalled)
}

func TestAllocAndFreeAreCalledCorrectlyPerItemFromPipe(t *testing.T) {
	sourceSlice := []type1{
		{
			Name: "John",
			Age:  32,
		},
		{
			Name: "Fred",
			Age:  55,
		},
		{
			Name: "Islam",
			Age:  35,
		},
	}
	source := enumerable.New(sourceSlice)

	runtime := newRuntime()

	wasmModule, err := engine.NewModule(runtime, modules.WasmPath4)
	if err != nil {
		t.Error(err)
	}

	wasmModuleInstance, err := engine.NewInstance(
		wasmModule,
		// We rename Name to Name, creating a transform that does nothing, but
		// will still affect the pipeline used, causing the native module to use
		// the `fromPipe` pipe, instead of the `fromSource` pipe.
		map[string]any{
			"src": "Name",
			"dst": "Name",
		},
	)
	if err != nil {
		t.Error(err)
	}

	timesAllocCalled := 0
	timesFreeCalled := 0
	var lastItemIndex module.MemSize
	var lastItemSize module.MemSize

	testModule := newNativeModule()
	results := engine.Append[type1, type2](
		source,
		wasmModuleInstance,
		module.Instance{
			Alloc: func(size module.MemSize) (module.MemSize, error) {
				lastItemIndex++
				timesAllocCalled++
				lastItemSize = size
				return lastItemIndex, nil
			},
			Free: func(ptr, size module.MemSize) error {
				timesFreeCalled++
				// Assert that we are freeing the last allocated item
				assert.Equal(t, ptr, lastItemIndex)
				assert.Equal(t, size, lastItemSize)
				return nil
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

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)

	// Assert alloc and free were called once per item
	assert.Equal(t, len(sourceSlice), timesAllocCalled)
	assert.Equal(t, len(sourceSlice), timesFreeCalled)
}

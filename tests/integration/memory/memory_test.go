// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package memory

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wasmtime"
	"github.com/lens-vm/lens/host-go/runtimes/wazero"
	"github.com/lens-vm/lens/tests/modules"
)

// The maximum allocation that can be made in one alloc call across all tested wasm langs.
//
// Language specific limits are as follows:
// - Rust: math.MaxInt32
// - AssemblyScript: math.MaxInt32 / 4
const maxAllocationSize int32 = math.MaxInt32 / 4

var runtimeConstructors = []func() module.Runtime{
	wasmtime.New,
	wazero.New,
	// wasmer doesn't work/compile on windows, so we append it to the set via not_windows_runtimes_test.go
	//wasmer.New,
}

// Note: Testing that lens modules do not leak memory when passing items to/from via the cli would be slow,
// for now we make do with testing that alloc and free do work, in all test modules and all runtimes.
// Indiviual host engines will need to write their own tests to ensure that alloc and free are called
// correctly.

func TestAllocErrorsWhenOutOfMemory(t *testing.T) {
	for _, modulePath := range modules.AllModules {
		for _, runtimeConstructor := range runtimeConstructors {
			runtime := runtimeConstructor()

			module, err := engine.NewModule(runtime, modulePath)
			require.NoError(t, err)

			instance, err := module.NewInstance("transform")
			require.NoError(t, err)

			indexes := make([]int32, 7)
			for i := 1; i <= 7; i++ {
				index, err := instance.Alloc(maxAllocationSize)
				require.NoError(t, err)

				indexes = append(indexes, index)
			}
			// The 8th call should fail, as we run out of 32bit memory
			_, err = instance.Alloc(maxAllocationSize)
			require.Error(t, err)

			for _, index := range indexes {
				// Free the allocations at the end of the runtime-loop, or
				// we'll have out of memory problems on smaller machines
				_ = instance.Free(index, maxAllocationSize)
			}
		}
	}
}

func TestAllocThenFreePastMemoryLimitDoesNotError(t *testing.T) {
	for _, modulePath := range modules.AllModules {
		for _, runtimeConstructor := range runtimeConstructors {
			runtime := runtimeConstructor()

			module, err := engine.NewModule(runtime, modulePath)
			require.NoError(t, err)

			instance, err := module.NewInstance("transform")
			require.NoError(t, err)

			// The 8th call would fail if free did not work, as we would
			// run out of 32bit memory.  Assuming free cannot partially
			// succeed - the testing for which would be too expensive
			// to bother right now.
			for i := 1; i <= 8; i++ {
				index, err := instance.Alloc(maxAllocationSize)
				require.NoError(t, err)
				err = instance.Free(index, maxAllocationSize)
				require.NoError(t, err)
			}
		}
	}
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !windows

package memory

import "github.com/lens-vm/lens/host-go/runtimes/wasmer"

func init() {
	// Wasmer doesn't work on windows so we add it to the set this way
	runtimeConstructors = append(runtimeConstructors, wasmer.New)
}

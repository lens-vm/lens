// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tests

import (
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wazero"
)

type type1 struct {
	Name string
	Age  int
}

type type2 struct {
	FullName string
	Age      int
}

func newRuntime() module.Runtime {
	return wazero.New()
}

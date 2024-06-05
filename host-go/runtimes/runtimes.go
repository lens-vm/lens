//go:build !js

package runtimes

import (
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wazero"
)

func Default() module.Runtime {
	return wazero.New()
}

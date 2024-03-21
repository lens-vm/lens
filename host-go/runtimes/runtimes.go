//go:build !js

package runtimes

import (
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wasmtime"
)

func Default() module.Runtime {
	return wasmtime.New()
}

package runtimes

import (
	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/js"
)

func Default() module.Runtime {
	return js.New()
}

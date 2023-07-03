package tests

import (
	"os"
	"testing"

	"github.com/lens-vm/lens/tests/modules"

	"github.com/wasmerio/wasmer-go/wasmer"
)

var (
	wasmModule *wasmer.Module
	store      *wasmer.Store
	instance   *wasmer.Instance
)

func BenchmarkWasmerRuntime(b *testing.B) {
	_, err := os.ReadFile(modules.WasmPath1)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		engine := wasmer.NewEngine()
		store = wasmer.NewStore(engine)
	}
}

func BenchmarkWasmerRuntimeWithModule(b *testing.B) {
	content, err := os.ReadFile(modules.WasmPath1)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		engine := wasmer.NewEngine()
		store := wasmer.NewStore(engine)

		wasmModule, err = wasmer.NewModule(store, content)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkWasmerRuntimeWithModuleAndInstance(b *testing.B) {
	content, err := os.ReadFile(modules.WasmPath1)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		engine := wasmer.NewEngine()
		store := wasmer.NewStore(engine)

		wasmModule, err = wasmer.NewModule(store, content)
		if err != nil {
			b.Error(err)
			return
		}

		importObject := wasmer.NewImportObject()
		instance, err = wasmer.NewInstance(wasmModule, importObject)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkWasmerRuntimeWithInstance(b *testing.B) {
	content, err := os.ReadFile(modules.WasmPath1)
	if err != nil {
		b.Error(err)
		return
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	wasmModule, err := wasmer.NewModule(store, content)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		importObject := wasmer.NewImportObject()
		instance, err = wasmer.NewInstance(wasmModule, importObject)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

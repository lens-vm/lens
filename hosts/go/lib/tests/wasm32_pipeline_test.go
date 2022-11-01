package tests

import (
	"lens-host/lib"
	"lens-host/lib/enumerable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWasm32PipelineFromSourceAsFull(t *testing.T) {
	module, err := lib.LoadModule(wasmPath1)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe := lib.AppendLens[type1, type2](source, module)

	hasNext, err := pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "John",
		Age:      32,
	}, pipe.Value())

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFull(t *testing.T) {
	module1, err := lib.LoadModule(wasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := lib.LoadModule(wasmPath2)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := lib.AppendLens[type1, type2](source, module1)
	pipe2 := lib.AppendLens[type2, type2](pipe1, module2)

	hasNext, err := pipe2.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "John",
		Age:      33,
	}, pipe2.Value())

	hasNext, err = pipe2.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFullToModuleAsFull(t *testing.T) {
	module1, err := lib.LoadModule(wasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := lib.LoadModule(wasmPath2)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := lib.AppendLens[type1, type2](source, module1)
	pipe2 := lib.AppendLens[type2, type2](pipe1, module2)
	pipe3 := lib.AppendLens[type2, type2](pipe2, module2)

	hasNext, err := pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "John",
		Age:      34,
	}, pipe3.Value())

	hasNext, err = pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

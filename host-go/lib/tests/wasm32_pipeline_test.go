package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/lib"
	"github.com/lens-vm/lens/host-go/lib/enumerable"

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

	pipe := lib.Append[type1, type2](source, module)

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

	pipe1 := lib.Append[type1, type2](source, module1)
	pipe2 := lib.Append[type2, type2](pipe1, module2)

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

	pipe1 := lib.Append[type1, type2](source, module1)
	pipe2 := lib.Append[type2, type2](pipe1, module2)
	pipe3 := lib.Append[type2, type2](pipe2, module2)

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

func TestWasm32PipelineFromSourceAsFullToModuleAsFullToASModuleAsFull(t *testing.T) {
	module1, err := lib.LoadModule(wasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := lib.LoadModule(wasmPath2)
	if err != nil {
		t.Error(err)
	}
	module3, err := lib.LoadModule(wasmPath3)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := lib.Append[type1, type2](source, module1)
	pipe2 := lib.Append[type2, type2](pipe1, module2)
	pipe3 := lib.Append[type2, type2](pipe2, module3)

	hasNext, err := pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "John",
		Age:      43,
	}, pipe3.Value())

	hasNext, err = pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

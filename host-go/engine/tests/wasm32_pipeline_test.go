package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/tests/modules"
	"github.com/sourcenetwork/immutable/enumerable"

	"github.com/stretchr/testify/assert"
)

func TestWasm32PipelineFromSourceAsFull(t *testing.T) {
	module, err := engine.LoadModule(modules.WasmPath1)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe := engine.Append[type1, type2](source, module)

	hasNext, err := pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe.Value()
	assert.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      32,
	}, val)

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFull(t *testing.T) {
	module1, err := engine.LoadModule(modules.WasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := engine.LoadModule(modules.WasmPath2)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := engine.Append[type1, type2](source, module1)
	pipe2 := engine.Append[type2, type2](pipe1, module2)

	hasNext, err := pipe2.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe2.Value()
	assert.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      33,
	}, val)

	hasNext, err = pipe2.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFullToModuleAsFull(t *testing.T) {
	module1, err := engine.LoadModule(modules.WasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := engine.LoadModule(modules.WasmPath2)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := engine.Append[type1, type2](source, module1)
	pipe2 := engine.Append[type2, type2](pipe1, module2)
	pipe3 := engine.Append[type2, type2](pipe2, module2)

	hasNext, err := pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe3.Value()
	assert.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      34,
	}, val)

	hasNext, err = pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFullToASModuleAsFull(t *testing.T) {
	module1, err := engine.LoadModule(modules.WasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := engine.LoadModule(modules.WasmPath2)
	if err != nil {
		t.Error(err)
	}
	module3, err := engine.LoadModule(modules.WasmPath3)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe1 := engine.Append[type1, type2](source, module1)
	pipe2 := engine.Append[type2, type2](pipe1, module2)
	pipe3 := engine.Append[type2, type2](pipe2, module3)

	hasNext, err := pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe3.Value()
	assert.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      43,
	}, val)

	hasNext, err = pipe3.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

func TestWasm32PipelineFromSourceAsFullToModuleAsFullToModuleAsFullWithSingleAppend(t *testing.T) {
	module1, err := engine.LoadModule(modules.WasmPath1)
	if err != nil {
		t.Error(err)
	}
	module2, err := engine.LoadModule(modules.WasmPath2)
	if err != nil {
		t.Error(err)
	}

	input := type1{
		Name: "John",
		Age:  32,
	}
	source := enumerable.New([]type1{input})

	pipe := engine.Append[type1, type2](
		source,
		module1,
		module2,
		module2,
	)

	hasNext, err := pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe.Value()
	assert.Nil(t, err)
	assert.Equal(t, type2{
		FullName: "John",
		Age:      34,
	}, val)

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

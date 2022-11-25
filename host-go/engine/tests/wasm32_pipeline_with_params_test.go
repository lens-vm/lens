package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/tests/modules"
	"github.com/sourcenetwork/immutable/enumerable"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWasm32PipelineWithAddtionalParams(t *testing.T) {
	module, err := engine.LoadModule(
		modules.WasmPath4,
		map[string]any{
			"src": "Name",
			"dst": "FullName",
		},
	)
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
	require.Nil(t, err)
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

func TestWasm32PipelineMultipleModulesAndWithAddtionalParams(t *testing.T) {
	module1, err := engine.LoadModule(
		modules.WasmPath4,
		map[string]any{
			"src": "Name",
			"dst": "FirstName",
		},
	)
	if err != nil {
		t.Error(err)
	}

	module2, err := engine.LoadModule(
		modules.WasmPath4,
		map[string]any{
			"src": "FirstName",
			"dst": "FullName",
		},
	)
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
	)

	hasNext, err := pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe.Value()
	require.Nil(t, err)
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

func TestWasm32PipelineWithAddtionalParamsErrors(t *testing.T) {
	module, err := engine.LoadModule(
		modules.WasmPath4,
		map[string]any{
			"src": "NotAField",
			"dst": "FullName",
		},
	)
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

	_, err = pipe.Value()
	assert.ErrorContains(t, err, "NotAField was not found")
}

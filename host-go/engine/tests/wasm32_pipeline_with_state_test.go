package tests

import (
	"testing"

	"github.com/lens-vm/lens/host-go/engine"
	"github.com/lens-vm/lens/tests/modules"
	"github.com/sourcenetwork/immutable/enumerable"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test asserts how state can be shared between pipe stages if the same module-instance is
// appended multiple times.
func TestWasm32PipelineWithSharedState(t *testing.T) {
	type Value struct {
		Id   int
		Name string
	}

	module, err := engine.LoadModule(modules.WasmPath5)
	if err != nil {
		t.Error(err)
	}

	source := enumerable.New([]Value{
		{
			Name: "John",
		},
		{
			Name: "Shahzad",
		},
		{
			Name: "Addo",
		},
	})

	pipe := engine.Append[Value, Value](source, module)
	pipe = engine.Append[Value, Value](pipe, module)
	pipe = engine.Append[Value, Value](pipe, module)

	hasNext, err := pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err := pipe.Value()
	require.Nil(t, err)
	assert.Equal(t, Value{
		Id:   3,
		Name: "John",
	}, val)

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err = pipe.Value()
	require.Nil(t, err)
	assert.Equal(t, Value{
		Id:   6,
		Name: "Shahzad",
	}, val)

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	val, err = pipe.Value()
	require.Nil(t, err)
	assert.Equal(t, Value{
		Id:   9,
		Name: "Addo",
	}, val)

	hasNext, err = pipe.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

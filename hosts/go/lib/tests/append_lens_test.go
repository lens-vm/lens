package tests

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"lens-host/lib"
	"lens-host/lib/enumerable"
	"lens-host/lib/module"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAppendLensWithoutWasm asserts that AppendLens can function independently of anything wasm related.
// This may be very important at a later date, as it means this function can (in-theory, minus any internal
// dependencies) can be executed in a wasm/wasi runtime. It should also help flag any changes to AppendLens'
// externally visible behaviour that may otherwise be hidden (e.g. if otherwise testing using wasm modules).
//
// TestAppendLensWithoutWasm also helps document how AppendLens works.
//
// It is not recomended to actually use AppendLens like this - if you are not transforming data via a wasm module,
// you should probably be using enumerable.Select or similar instead.
func TestAppendLensWithoutWasm(t *testing.T) {
	sourceSlice := []type1{
		{
			Name: "John",
			Age:  32,
		},
		{
			Name: "Fred",
			Age:  55,
		},
	}
	source := enumerable.New(sourceSlice)

	memory := make([]byte, math.MaxUint16)
	results := lib.AppendLens[type1, type2](
		source,
		module.Module{
			Alloc: func(size module.MemSize) (module.MemSize, error) {
				var arbitraryIndex module.MemSize = 5
				return arbitraryIndex, nil
			},
			Transform: func(startIndex module.MemSize, additionalParams ...any) (module.MemSize, error) {
				resultBuffer := make([]byte, module.LenSize)
				copy(resultBuffer, memory[startIndex:startIndex+module.LenSize])
				var inputLen module.LenType
				buf := bytes.NewReader(resultBuffer)
				_ = binary.Read(buf, module.LenByteOrder, &inputLen)

				sourceJson := memory[startIndex+module.LenSize : startIndex+module.MemSize(inputLen)+module.LenSize]
				sourceItem := &type1{}

				err := json.Unmarshal([]byte(sourceJson), sourceItem)
				if err != nil {
					return 0, err
				}

				resultItem := type2{
					FullName: sourceItem.Name,
					Age:      sourceItem.Age,
				}

				returnBytes, err := json.Marshal(resultItem)
				if err != nil {
					return 0, err
				}

				returnLen := module.LenType(len(returnBytes))
				writer := bytes.NewBuffer([]byte{})
				err = binary.Write(writer, module.LenByteOrder, returnLen)
				if err != nil {
					return 0, err
				}

				arbitraryReturnIndex := math.MaxUint16 / 2
				dst := memory[arbitraryReturnIndex:]
				copy(dst, writer.Bytes())
				copy(dst[module.LenSize:], returnBytes)

				return module.MemSize(arbitraryReturnIndex), nil
			},
			GetData: func() []byte {
				return memory
			},
		},
	)

	hasNext, err := results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "John",
		Age:      32,
	}, results.Value())

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, hasNext)

	assert.Equal(t, type2{
		FullName: "Fred",
		Age:      55,
	}, results.Value())

	hasNext, err = results.Next()
	if err != nil {
		t.Error(err)
	}
	assert.False(t, hasNext)
}

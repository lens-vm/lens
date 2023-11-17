// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tests

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/lens-vm/lens/host-go/runtimes/wasmtime"
)

type type1 struct {
	Name string
	Age  int
}

type type2 struct {
	FullName string
	Age      int
}

func newRuntime() module.Runtime {
	return wasmtime.New()
}

// nativeModule is a native Go lens module that transforms `type1` items to `type2`.
//
// It is used for testing puposes only.
type nativeModule struct {
	memory []byte
}

func newNativeModule() *nativeModule {
	return &nativeModule{
		memory: make([]byte, math.MaxUint16),
	}
}

func (m *nativeModule) nativeTransform(startIndex module.MemSize) (module.MemSize, error) {
	typeBuffer := make([]byte, module.TypeIdSize)
	copy(typeBuffer, m.memory[startIndex:startIndex+module.TypeIdSize])
	var inputTypeId module.TypeIdType
	buf := bytes.NewReader(typeBuffer)
	err := binary.Read(buf, module.TypeIdByteOrder, &inputTypeId)
	if err != nil {
		return 0, err
	}

	lenBuffer := make([]byte, module.LenSize)
	copy(lenBuffer, m.memory[startIndex+module.TypeIdSize:startIndex+module.TypeIdSize+module.LenSize])
	var inputLen module.LenType
	buf = bytes.NewReader(lenBuffer)
	err = binary.Read(buf, module.LenByteOrder, &inputLen)
	if err != nil {
		return 0, err
	}

	sourceJson := m.memory[startIndex+module.TypeIdSize+module.LenSize : startIndex+module.TypeIdSize+module.MemSize(inputLen)+module.LenSize]
	sourceItem := &type1{}

	err = json.Unmarshal([]byte(sourceJson), sourceItem)
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

	typeWriter := bytes.NewBuffer([]byte{})
	err = binary.Write(typeWriter, module.TypeIdByteOrder, int8(1))
	if err != nil {
		return 0, err
	}

	returnLen := module.LenType(len(returnBytes))
	lenWriter := bytes.NewBuffer([]byte{})
	err = binary.Write(lenWriter, module.LenByteOrder, returnLen)
	if err != nil {
		return 0, err
	}

	arbitraryReturnIndex := math.MaxUint16 / 2
	dst := m.memory[arbitraryReturnIndex:]
	copy(dst, typeWriter.Bytes())
	copy(dst[module.TypeIdSize:], lenWriter.Bytes())
	copy(dst[module.TypeIdSize+module.LenSize:], returnBytes)

	return module.MemSize(arbitraryReturnIndex), nil
}

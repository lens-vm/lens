// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/lens-vm/lens/host-go/engine/module"
)

// GetItem returns the item at the given index.  This includes the length specifier.
func GetItem(src []byte, startIndex module.MemSize) ([]byte, error) {
	typeBuffer := make([]byte, module.TypeIdSize)
	copy(typeBuffer, src[startIndex:startIndex+module.TypeIdSize])
	var typeId module.TypeIdType
	reader := bytes.NewReader(typeBuffer)
	err := binary.Read(reader, module.TypeIdByteOrder, &typeId)
	if err != nil {
		return nil, err
	}

	if typeId == module.NilTypeID {
		return nil, nil
	}

	lenBuffer := make([]byte, module.LenSize)
	copy(lenBuffer, src[startIndex+module.TypeIdSize:startIndex+module.TypeIdSize+module.LenSize])
	var len module.LenType
	reader = bytes.NewReader(lenBuffer)
	err = binary.Read(reader, module.LenByteOrder, &len)
	if err != nil {
		return nil, err
	}

	if module.IsError(typeId) {
		return nil, errors.New(
			string(
				src[startIndex+module.TypeIdSize+module.LenSize : startIndex+module.TypeIdSize+module.MemSize(len)+module.LenSize],
			),
		)
	}

	// todo - the end index of this is untested, as it will only affect performance atm if it is longer than desired
	// unless it overwrites adjacent stuff
	return src[startIndex : startIndex+module.TypeIdSize+module.MemSize(len)+module.LenSize], nil
}

// WriteItem calculates the length specifier for the given source object and then writes both specifier
// and item to the destination.
func WriteItem(typeId module.TypeIdType, src []byte, dst []byte) error {
	typeWriter := bytes.NewBuffer([]byte{})
	err := binary.Write(typeWriter, module.TypeIdByteOrder, typeId)
	if err != nil {
		return err
	}

	len := module.LenType(len(src))
	lenWriter := bytes.NewBuffer([]byte{})
	err = binary.Write(lenWriter, module.LenByteOrder, len)
	if err != nil {
		return err
	}

	copy(dst, typeWriter.Bytes())
	copy(dst[module.TypeIdSize:], lenWriter.Bytes())
	copy(dst[module.TypeIdSize+module.LenSize:], src)

	return nil
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/lens-vm/lens/host-go/engine/module"
)

// ReadTypeId returns the type id of the next item from the given reader.
func ReadTypeId(r io.Reader) (module.TypeIdType, error) {
	typeBuffer := make([]byte, module.TypeIdSize)
	typeReader := bytes.NewReader(typeBuffer)

	_, err := r.Read(typeBuffer)
	if err != nil {
		return 0, err
	}
	var typeId module.TypeIdType
	err = binary.Read(typeReader, module.TypeIdByteOrder, &typeId)
	if err != nil {
		return 0, err
	}
	return typeId, nil
}

// ReadItem returns the bytes of the next item from the given reader.
func ReadItem(r io.Reader) ([]byte, error) {
	typeId, err := ReadTypeId(r)
	if err != nil {
		return nil, err
	}

	// type is nil so nothing else to read
	if typeId == module.NilTypeID {
		return nil, nil
	}

	lenBuffer := make([]byte, module.LenSize)
	lenReader := bytes.NewReader(lenBuffer)

	// read the item length
	_, err = r.Read(lenBuffer)
	if err != nil {
		return nil, err
	}
	var len module.LenType
	err = binary.Read(lenReader, module.LenByteOrder, &len)
	if err != nil {
		return nil, err
	}

	// read the item bytes
	data := make([]byte, len)
	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}

	if typeId.IsError() {
		return nil, errors.New(string(data))
	}
	return data, nil
}

func WriteItem(w io.Writer, id module.TypeIdType, data []byte) error {
	// write the item type id
	err := binary.Write(w, module.TypeIdByteOrder, id)
	if err != nil {
		return err
	}

	// end of stream messages have no value component that needs writing
	if id.IsEOS() {
		return nil
	}

	// write the item length
	err = binary.Write(w, module.LenByteOrder, module.LenType(len(data)))
	if err != nil {
		return err
	}

	// write the item bytes
	_, err = w.Write(data)
	return err
}

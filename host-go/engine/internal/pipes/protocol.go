package pipes

import (
	"bytes"
	"encoding/binary"

	"github.com/lens-vm/lens/host-go/engine/module"
)

// getItem returns the item at the given index.  This includes the length specifier.
func getItem(src []byte, startIndex module.MemSize) ([]byte, error) {
	lenBuffer := make([]byte, module.LenSize)
	copy(lenBuffer, src[startIndex:startIndex+module.LenSize])
	var len module.LenType
	reader := bytes.NewReader(lenBuffer)
	err := binary.Read(reader, module.LenByteOrder, &len)
	if err != nil {
		return nil, err
	}

	// todo - the end index of this is untested, as it will only affect performance atm if it is longer than desired
	// unless it overwrites adjacent stuff
	return src[startIndex : startIndex+module.MemSize(len)+module.LenSize], nil
}

// WriteItem calculates the length specifier for the given source object and then writes both specifier
// and item to the destination.
func WriteItem(src []byte, dst []byte) error {
	len := module.LenType(len(src))
	writer := bytes.NewBuffer([]byte{})
	err := binary.Write(writer, module.LenByteOrder, len)
	if err != nil {
		return err
	}

	copy(dst, writer.Bytes())
	copy(dst[module.LenSize:], src)

	return nil
}

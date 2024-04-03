// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/sourcenetwork/immutable/enumerable"
)

type fromSource[TSource any, TResult any] struct {
	source   enumerable.Enumerable[TSource]
	instance module.Instance

	currentIndex module.MemSize
}

func NewFromSource[TSource any, TResult any](
	source enumerable.Enumerable[TSource],
	instance module.Instance,
) Pipe[TResult] {
	return &fromSource[TSource, TResult]{
		source:   source,
		instance: instance,
	}
}

var _ Pipe[int] = (*fromSource[bool, int])(nil)

func (s *fromSource[TSource, TResult]) Next() (bool, error) {
	index, err := s.instance.Transform(s.mustGetNext)
	if err != nil {
		return false, err
	}
	typeId, err := ReadTypeId(s.instance.Memory(index))
	if err != nil {
		return false, err
	}
	if typeId.IsEOS() {
		return false, nil
	}

	s.currentIndex = index
	return true, nil
}

func (s *fromSource[TSource, TResult]) Value() (TResult, error) {
	var result TResult

	id, data, err := ReadItem(s.instance.Memory(s.currentIndex))
	if err != nil {
		return result, err
	}
	if id.IsError() {
		return result, errors.New(string(data))
	}
	if id != module.JSONTypeID {
		return result, nil
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (s *fromSource[TSource, TResult]) Bytes() ([]byte, error) {
	id, data, err := ReadItem(s.instance.Memory(s.currentIndex))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := WriteItem(&out, id, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (s *fromSource[TSource, TResult]) Reset() {
	s.source.Reset()
}

// mustGetNext tries to get the next value from source and copy it into the memory buffer.
//
// If there are no more values in source it will write an EOS message. If an error is
// generated, it will attempt to write the error to the memory buffer - if the writing of the
// error to the buffer fails it will panic.
func (s *fromSource[TSource, TResult]) mustGetNext() module.MemSize {
	index, err := s.getNext()
	if err != nil {
		return mustWriteErr(s.instance, err)
	}

	return index
}

func (s *fromSource[TSource, TResult]) getNext() (module.MemSize, error) {
	hasNext, err := s.source.Next()
	if err != nil {
		return 0, err
	}
	if !hasNext {
		return writeEOS(s.instance)
	}

	sourceItem, err := s.source.Value()
	if err != nil {
		return 0, err
	}
	value, err := json.Marshal(sourceItem)
	if err != nil {
		return 0, err
	}
	// allocate space for the next item
	index, err := s.instance.Alloc(module.TypeIdSize + module.LenSize + module.MemSize(len(value)))
	if err != nil {
		return 0, err
	}
	// write the item to memory
	err = WriteItem(s.instance.Memory(index), module.JSONTypeID, value)
	if err != nil {
		return 0, err
	}
	return index, nil
}

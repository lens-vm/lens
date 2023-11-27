// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import (
	"encoding/json"

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

	if module.TypeIdType(s.instance.GetData()[index]).IsEOS() {
		return false, nil
	}

	s.currentIndex = index
	return true, nil
}

func (s *fromSource[TSource, TResult]) Value() (TResult, error) {
	var t TResult

	item, err := GetItem(s.instance.GetData(), s.currentIndex)
	if err != nil || item == nil {
		return t, err
	}
	jsonBytes := item[module.TypeIdSize+module.LenSize:]

	result := &t
	err = json.Unmarshal(jsonBytes, result)
	if err != nil {
		return t, err
	}

	return *result, nil
}

func (s *fromSource[TSource, TResult]) Bytes() ([]byte, error) {
	return GetItem(s.instance.GetData(), s.currentIndex)
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

	index, err := s.instance.Alloc(module.TypeIdSize + module.LenSize + module.MemSize(len(value)))
	if err != nil {
		return 0, err
	}

	err = WriteItem(module.JSONTypeID, value, s.instance.GetData()[index:])
	if err != nil {
		return 0, err
	}
	return index, nil
}

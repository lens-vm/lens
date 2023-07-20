// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
)

type fromPipe[TSource any, TResult any] struct {
	source   Pipe[TSource]
	instance module.Instance

	currentIndex module.MemSize
}

func NewFromPipe[TSource any, TResult any](
	source Pipe[TSource],
	instance module.Instance,
) Pipe[TResult] {
	return &fromPipe[TSource, TResult]{
		source:   source,
		instance: instance,
	}
}

var _ Pipe[int] = (*fromPipe[bool, int])(nil)

func (s *fromPipe[TSource, TResult]) Next() (bool, error) {
	hasNext, err := s.source.Next()
	if !hasNext || err != nil {
		return hasNext, err
	}

	value, err := s.source.Bytes()
	if err != nil {
		return false, nil
	}

	// We do this here to keep the work (and errors) in the `Next` call
	result, err := s.transport(value)
	if err != nil {
		return false, nil
	}

	s.currentIndex = result
	return true, nil
}

func (s *fromPipe[TSource, TResult]) Value() (TResult, error) {
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

func (s *fromPipe[TSource, TResult]) Bytes() ([]byte, error) {
	return GetItem(s.instance.GetData(), s.currentIndex)
}

func (s *fromPipe[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromPipe[TSource, TResult]) transport(sourceItem []byte) (module.MemSize, error) {
	index, err := s.instance.Alloc(module.MemSize(len(sourceItem)))
	if err != nil {
		return 0, err
	}

	copy(s.instance.GetData()[index:], sourceItem)

	index, err = s.instance.Transform(index)
	if err != nil {
		return 0, err
	}

	return index, nil
}
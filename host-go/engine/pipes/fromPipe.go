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

func (p *fromPipe[TSource, TResult]) Next() (bool, error) {
	index, err := p.instance.Transform(p.mustGetNext)
	if err != nil {
		return false, err
	}
	typeId, err := ReadTypeId(p.instance.Memory(index))
	if err != nil {
		return false, err
	}
	if typeId.IsEOS() {
		return false, nil
	}

	p.currentIndex = index
	return true, nil
}

func (p *fromPipe[TSource, TResult]) Value() (TResult, error) {
	var t TResult

	item, err := ReadItem(p.instance.Memory(p.currentIndex))
	if err != nil || item == nil {
		return t, err
	}

	result := &t
	err = json.Unmarshal(item, result)
	if err != nil {
		return t, err
	}
	return *result, nil
}

func (p *fromPipe[TSource, TResult]) Bytes() ([]byte, error) {
	return ReadItem(p.instance.Memory(p.currentIndex))
}

func (p *fromPipe[TSource, TResult]) Reset() {
	p.source.Reset()
}

// mustGetNext tries to get the next value from source and copy it into the memory buffer.
//
// If there are no more values in source it will write an EOS message. If an error is
// generated, it will attempt to write the error to the memory buffer - if the writing of the
// error to the buffer fails it will panic.
func (p *fromPipe[TSource, TResult]) mustGetNext() module.MemSize {
	index, err := p.getNext()
	if err == nil {
		return index
	}
	data := []byte(err.Error())
	// allocate space for the error
	index, err = p.instance.Alloc(module.TypeIdSize + module.LenSize + int32(len(data)))
	if err != nil {
		panic(err)
	}
	// write the error starting at the allocated index
	err = WriteItem(p.instance.Memory(index), module.ErrTypeID, data)
	if err != nil {
		panic(err)
	}
	return index
}

func (p *fromPipe[TSource, TResult]) getNext() (module.MemSize, error) {
	hasNext, err := p.source.Next()
	if err != nil {
		return 0, err
	}
	if !hasNext {
		// allocate space for EOS
		index, err := p.instance.Alloc(module.TypeIdSize)
		if err != nil {
			return 0, err
		}
		// write EOS to memory
		err = WriteItem(p.instance.Memory(index), module.EOSTypeID, nil)
		if err != nil {
			return 0, err
		}
		return index, nil
	}

	value, err := p.source.Bytes()
	if err != nil {
		return 0, err
	}
	// allocate space for the next item
	index, err := p.instance.Alloc(module.TypeIdSize + module.LenSize + module.MemSize(len(value)))
	if err != nil {
		return 0, err
	}
	// write the item to memory
	err = WriteItem(p.instance.Memory(index), module.JSONTypeID, value)
	if err != nil {
		return 0, err
	}
	return index, nil
}

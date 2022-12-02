package pipes

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
)

type fromPipe[TSource any, TResult any] struct {
	source Pipe[TSource]
	module module.Module

	currentIndex module.MemSize
}

func NewFromPipe[TSource any, TResult any](
	source Pipe[TSource],
	module module.Module,
) Pipe[TResult] {
	return &fromPipe[TSource, TResult]{
		source: source,
		module: module,
	}
}

var _ Pipe[int] = (*fromPipe[bool, int])(nil)

func (s *fromPipe[TSource, TResult]) Next() (bool, error) {
	index, err := s.module.Transform(s.mustGetNext)
	if err != nil {
		return false, err
	}

	if module.IsEOS(module.TypeIdType(s.module.GetData()[index])) {
		return false, nil
	}

	s.currentIndex = index
	return true, nil
}

func (s *fromPipe[TSource, TResult]) Value() (TResult, error) {
	var t TResult

	item, err := GetItem(s.module.GetData(), s.currentIndex)
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
	return GetItem(s.module.GetData(), s.currentIndex)
}

func (s *fromPipe[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromPipe[TSource, TResult]) mustGetNext() module.MemSize {
	index, err := s.getNext()
	if err != nil {
		return mustWriteErr(s.module, err)
	}

	return index
}

func (s *fromPipe[TSource, TResult]) getNext() (module.MemSize, error) {
	hasNext, err := s.source.Next()
	if err != nil {
		return 0, err
	}

	if !hasNext {
		return writeEOS(s.module)
	}

	value, err := s.source.Bytes()
	if err != nil {
		return 0, err
	}

	index, err := s.module.Alloc(module.MemSize(len(value)))
	if err != nil {
		return 0, err
	}

	copy(s.module.GetData()[index:], value)
	return index, nil
}

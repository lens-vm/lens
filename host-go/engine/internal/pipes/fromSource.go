package pipes

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
	"github.com/sourcenetwork/immutable/enumerable"
)

type fromSource[TSource any, TResult any] struct {
	source enumerable.Enumerable[TSource]
	module module.Module

	currentIndex module.MemSize
}

func NewFromSource[TSource any, TResult any](
	source enumerable.Enumerable[TSource],
	module module.Module,
) Pipe[TResult] {
	return &fromSource[TSource, TResult]{
		source: source,
		module: module,
	}
}

var _ Pipe[int] = (*fromSource[bool, int])(nil)

func (s *fromSource[TSource, TResult]) Next() (bool, error) {
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

func (s *fromSource[TSource, TResult]) Value() (TResult, error) {
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

func (s *fromSource[TSource, TResult]) Bytes() ([]byte, error) {
	return GetItem(s.module.GetData(), s.currentIndex)
}

func (s *fromSource[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromSource[TSource, TResult]) mustGetNext() module.MemSize {
	index, err := s.getNext()
	if err != nil {
		return mustWriteErr(s.module, err)
	}

	return index
}

func (s *fromSource[TSource, TResult]) getNext() (module.MemSize, error) {
	hasNext, err := s.source.Next()
	if err != nil {
		return 0, err
	}

	if !hasNext {
		return writeEOS(s.module)
	}

	sourceItem, err := s.source.Value()
	if err != nil {
		return 0, err
	}

	value, err := json.Marshal(sourceItem)
	if err != nil {
		return 0, err
	}

	index, err := s.module.Alloc(module.TypeIdSize + module.LenSize + module.MemSize(len(value)))
	if err != nil {
		return 0, err
	}

	err = WriteItem(module.JSONTypeID, value, s.module.GetData()[index:])
	if err != nil {
		return 0, err
	}
	return index, nil
}

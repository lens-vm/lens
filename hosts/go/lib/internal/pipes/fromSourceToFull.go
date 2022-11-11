package pipes

import (
	"encoding/json"
	"lens-host/lib/enumerable"
	"lens-host/lib/module"
)

type fromSourceToFull[TSource any, TResult any] struct {
	source enumerable.Enumerable[TSource]
	module module.Module

	currentIndex module.MemSize
}

func FromSourceToFull[TSource any, TResult any](
	source enumerable.Enumerable[TSource],
	module module.Module,
) Pipe[TResult] {
	return &fromSourceToFull[TSource, TResult]{
		source: source,
		module: module,
	}
}

var _ Pipe[int] = (*fromSourceToFull[bool, int])(nil)

func (s *fromSourceToFull[TSource, TResult]) Next() (bool, error) {
	hasNext, err := s.source.Next()
	if !hasNext || err != nil {
		return hasNext, err
	}

	value := s.source.Value()
	// We do this here to keep the work (and errors) in the `Next` call
	result, err := s.transport(value)
	if err != nil {
		return false, nil
	}

	s.currentIndex = result
	return true, nil
}

func (s *fromSourceToFull[TSource, TResult]) Value() TResult {
	item := getItem(s.module.GetData(), s.currentIndex)
	jsonStr := string(item[module.LenSize:])

	var t TResult
	result := &t
	err := json.Unmarshal([]byte(jsonStr), result)
	if err != nil {
		panic(err)
	}
	return *result
}

func (s *fromSourceToFull[TSource, TResult]) Bytes() []byte {
	return getItem(s.module.GetData(), s.currentIndex)
}

func (s *fromSourceToFull[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromSourceToFull[TSource, TResult]) transport(sourceItem TSource) (module.MemSize, error) {
	sourceBytes, err := json.Marshal(sourceItem)
	if err != nil {
		return 0, err
	}

	index, err := s.module.Alloc(module.MemSize(len(sourceBytes)) + module.LenSize)
	if err != nil {
		return 0, err
	}

	err = writeItem(sourceBytes, s.module.GetData()[index:])
	if err != nil {
		return 0, err
	}

	index, err = s.module.Transform(index)
	if err != nil {
		return 0, err
	}

	return index, nil
}

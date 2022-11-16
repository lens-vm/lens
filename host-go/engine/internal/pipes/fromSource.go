package pipes

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/enumerable"
	"github.com/lens-vm/lens/host-go/engine/module"
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

func (s *fromSource[TSource, TResult]) Value() TResult {
	item := getItem(s.module.GetData(), s.currentIndex)
	jsonStr := string(item[module.LenSize:])

	var t TResult
	result := &t
	err := json.Unmarshal([]byte(jsonStr), result)
	if err != nil {
		// TODO: We should return this instead of panicing
		// https://github.com/sourcenetwork/lens/issues/10
		panic(err)
	}

	return *result
}

func (s *fromSource[TSource, TResult]) Bytes() []byte {
	return getItem(s.module.GetData(), s.currentIndex)
}

func (s *fromSource[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromSource[TSource, TResult]) transport(sourceItem TSource) (module.MemSize, error) {
	sourceBytes, err := json.Marshal(sourceItem)
	if err != nil {
		return 0, err
	}

	index, err := s.module.Alloc(module.MemSize(len(sourceBytes)) + module.LenSize)
	if err != nil {
		return 0, err
	}

	err = WriteItem(sourceBytes, s.module.GetData()[index:])
	if err != nil {
		return 0, err
	}

	index, err = s.module.Transform(index)
	if err != nil {
		return 0, err
	}

	return index, nil
}

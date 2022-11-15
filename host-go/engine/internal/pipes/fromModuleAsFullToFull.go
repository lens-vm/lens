package pipes

import (
	"encoding/json"

	"github.com/lens-vm/lens/host-go/engine/module"
)

type fromModuleAsFullToFull[TSource any, TResult any] struct {
	source Pipe[TSource]
	module module.Module

	currentIndex module.MemSize
}

func FromModuleAsFullToFull[TSource any, TResult any](
	source Pipe[TSource],
	module module.Module,
) Pipe[TResult] {
	return &fromModuleAsFullToFull[TSource, TResult]{
		source: source,
		module: module,
	}
}

var _ Pipe[int] = (*fromModuleAsFullToFull[bool, int])(nil)

func (s *fromModuleAsFullToFull[TSource, TResult]) Next() (bool, error) {
	hasNext, err := s.source.Next()
	if !hasNext || err != nil {
		return hasNext, err
	}

	value := s.source.Bytes()
	// We do this here to keep the work (and errors) in the `Next` call
	result, err := s.transport(value)
	if err != nil {
		return false, nil
	}

	s.currentIndex = result
	return true, nil
}

func (s *fromModuleAsFullToFull[TSource, TResult]) Value() TResult {
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

func (s *fromModuleAsFullToFull[TSource, TResult]) Bytes() []byte {
	return getItem(s.module.GetData(), s.currentIndex)
}

func (s *fromModuleAsFullToFull[TSource, TResult]) Reset() {
	s.source.Reset()
}

func (s *fromModuleAsFullToFull[TSource, TResult]) transport(sourceItem []byte) (module.MemSize, error) {
	index, err := s.module.Alloc(module.MemSize(len(sourceItem)))
	if err != nil {
		return 0, err
	}

	copy(s.module.GetData()[index:], sourceItem)

	index, err = s.module.Transform(index)
	if err != nil {
		return 0, err
	}

	return index, nil
}

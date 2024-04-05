package wazero

import (
	"github.com/tetratelabs/wazero/api"
)

type memory struct {
	memory api.Memory
}

func newMemory(mem api.Memory) *memory {
	return &memory{mem}
}

func (m *memory) ReadAt(dst []byte, offset int64) (int, error) {
	out, _ := m.memory.Read(uint32(offset), uint32(len(dst)))
	n := copy(dst, out)
	return n, nil
}

func (m *memory) WriteAt(src []byte, offset int64) (int, error) {
	m.memory.Write(uint32(offset), src)
	return len(src), nil
}

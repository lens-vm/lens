package wazero

import (
	"github.com/tetratelabs/wazero/api"
)

type memory struct {
	memory api.Memory
	offset int32
}

func newMemory(mem api.Memory, offset int32) *memory {
	return &memory{memory: mem, offset: offset}
}

func (m *memory) Read(dst []byte) (int, error) {
	out, _ := m.memory.Read(uint32(m.offset), uint32(len(dst)))
	n := copy(dst, out)
	m.offset = m.offset + int32(n)
	return n, nil
}

func (m *memory) Write(src []byte) (int, error) {
	m.memory.Write(uint32(m.offset), src)
	m.offset = m.offset + int32(len(src))
	return len(src), nil
}

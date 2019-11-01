// Package memory implements concrete memory types
// along with a generic memory struct that can contain
// a number of memory spaces.
package memory

import "github.com/lucactt/gameboy/util/errors"

// Memory represents a general purpose memory from which bytes
// can be read and written.
//
// When an address is outside the memory, an error will be returned.
type Memory interface {
	GetByte(addr uint16) (byte, error)
	SetByte(addr uint16, value byte) error
	Accepts(addr uint16) bool
}

// space is a wrapper for a memory, that allows associating
// a start address with the start of the memory.
type space struct {
	start uint16
	mem   Memory
}

func (s *space) GetByte(addr uint16) (byte, error) {
	return s.mem.GetByte(addr - s.start)
}

func (s *space) SetByte(addr uint16, value byte) error {
	return s.mem.SetByte(addr-s.start, value)
}

func (s *space) Accepts(addr uint16) bool {
	return s.mem.Accepts(addr - s.start)
}

// MMU represents a Memory Management Unit that wraps many
// memories.
//
// It itself implements the Memory interface.
type MMU struct {
	spaces []*space
}

// GetByte returns the byte at the given address.
// If the address is outside every wrapped memory,
// it will return an error.
func (m *MMU) GetByte(addr uint16) (byte, error) {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			res, err := s.GetByte(addr)
			if err != nil {
				return 0, errors.E("get byte from memory failed", err, errors.Memory)
			}
			return res, nil
		}
	}

	return 0, errors.E("no memory space accepts addr", errors.CodeOutOfRange, errors.Memory)
}

// SetByte sets the byte at the given address.
// If the address is outside every wrapped memory,
// it will return an error.
func (m *MMU) SetByte(addr uint16, value byte) error {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			err := s.SetByte(addr, value)
			if err != nil {
				return errors.E("set byte to memory failed", err, errors.Memory)
			}
			return nil
		}
	}

	return errors.E("no memory space accepts addr", errors.CodeOutOfRange, errors.Memory)
}

// Accepts checks if any of the underlying memories
// accept the given address.
func (m *MMU) Accepts(addr uint16) bool {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			return true
		}
	}
	return false
}

// AddMem adds a memory to the MMU.
func (m *MMU) AddMem(start uint16, mem Memory) {
	m.spaces = append(m.spaces, &space{start, mem})
}

// Package memory implements concrete memory types
// along with a generic memory struct that can contain
// a number of memory spaces.
package memory

import "github.com/lucactt/gameboy/util/errors"

// Space represents a memory space from which bytes
// can be read and written.
//
// By convention, the start of the space is inclusive,
// while the end is non-inclusive.
// For example, a space between 0x0001 and 0x1000 will have
// addresses from 0x0001 and 0x0FFF.
//
// Also by convention, when an address is outside the space,
// an error will be returned.
type Space interface {
	GetByte(addr uint16) (byte, error)
	SetByte(addr uint16, value byte) error
	Accepts(addr uint16) bool
}

// Memory represents a generic memory that wraps many
// memory spaces.
// It itself implements the Space interface.
type Memory struct {
	spaces []Space
}

// GetByte returns the byte at the given address.
// If the address points to a non-existing memory space,
// an error will be returned.
func (m *Memory) GetByte(addr uint16) (byte, error) {
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
// If the address points to a non-existing memory space,
// an error will be returned.
func (m *Memory) SetByte(addr uint16, value byte) error {
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

// Accepts checks if any of the underlying memory spaces
// accept the given address.
func (m *Memory) Accepts(addr uint16) bool {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			return true
		}
	}
	return false
}

// AddSpace adds a memory space to the memory.
func (m *Memory) AddSpace(space Space) {
	m.spaces = append(m.spaces, space)
}

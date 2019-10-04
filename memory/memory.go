package memory

import "errors"

// AddressSpace represents a generic memory area that
// can be read and written to.
type AddressSpace interface {
	GetByte(addr uint16) (byte, error)
	SetByte(addr uint16, value byte) error
	Accepts(addr uint16) bool
}

// Memory represents the GameBoy memory.
// It implements the AddressSpace interface.
type Memory struct {
	spaces []AddressSpace
}

// GetByte returns the byte at the given address.
// If the address points to a non-existing memory area,
// an error will be returned.
func (m *Memory) GetByte(addr uint16) (byte, error) {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			res, err := s.GetByte(addr)
			if err != nil {
				return 0, err
			}
			return res, nil
		}
	}

	return 0, errors.New("")
}

// SetByte sets the value of the byte at the given address.
// If the address points to a non-existing memory area,
// an error will be returned.
func (m *Memory) SetByte(addr uint16, value byte) error {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			err := s.SetByte(addr, value)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("")
}

// AddSpace adds a memory bank to the memory, from which bytes
// can be read and written to.
// Will panic if the space addresses overlap with those of another space
// in the memory.
func (m *Memory) AddSpace(space AddressSpace) {
	m.spaces = append(m.spaces, space)
}

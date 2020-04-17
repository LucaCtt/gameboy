// Package mem implements various memory types.
package mem

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// Mem represents a general purpose memory from which bytes
// can be read and written.
//
// The addresses of a mem must always start at 0x00.
//
// When an address is not within the memory addresses range,
// an error must be returned.
//
// While not necessary, any code that uses this interface should
// check that Accepts(addr) returns true before getting or setting that address.
// If Accepts(addr) returns true, but getting or setting that address returns an exception,
// the calling code should panic, as that indicates a dev error (please report it).
type Mem interface {
	// GetByte returns the byte at the given address, or an error if
	// the mem doesn't accept the address.
	GetByte(addr uint16) (byte, error)

	// SetByte sets the byte at the given address to the given value, and
	// returns an error if the mem doesn't accept the address.
	SetByte(addr uint16, value byte) error

	// Accepts checks that the address is non-negative,
	// and strictly less than the memory length.
	Accepts(addr uint16) bool
}

// space is a wrapper for a mem, that associates
// a start address with the start of the memory.
//
// This is required because the interface Mem doesn't store a start addr.
type space struct {
	start uint16
	mem   Mem
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
// memories. Externally it behaves just like any memory.
//
// It implements the Mem interface.
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
				// If this happens it's because of a development error, so panic is ok.
				panic(errors.E(fmt.Sprintf("mem accepts %d, but GetByte returned error", addr), err))
			}
			return res, nil
		}
	}

	return 0, errors.E("no memory space accepts addr", errors.CodeOutOfRange, errors.Mem)
}

// SetByte sets the byte at the given address.
// If the address is outside every wrapped memory,
// it will return an error.
func (m *MMU) SetByte(addr uint16, value byte) error {
	for _, s := range m.spaces {
		if s.Accepts(addr) {
			err := s.SetByte(addr, value)
			if err != nil {
				panic(errors.E(fmt.Sprintf("mem accepts %d, but SetByte returned error", addr), err))
			}
			return nil
		}
	}

	return errors.E("no memory space accepts addr", errors.CodeOutOfRange, errors.Mem)
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

// AddMem adds a memory to the MMU at the given address.
//
// Note that if the MMU already contains a memory that "covers"
// (even some of) the addresses of the Mem to add, that memory
// will be the one to handle those addresses.
func (m *MMU) AddMem(start uint16, mem Mem) {
	m.spaces = append(m.spaces, &space{start, mem})
}

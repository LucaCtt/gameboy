package mem

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// ROM represents a generic Read Only Memory.
type ROM struct {
	rom []byte
}

// NewROM creates a new ROM with the given byte slice.
func NewROM(content []byte) *ROM {
	return &ROM{content}
}

// GetByte returns the byte at the given address.
// If the address is outside the memory, an
// error will be returned.
func (r *ROM) GetByte(addr uint16) (byte, error) {
	if !r.Accepts(addr) {
		return 0, errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return r.rom[addr], nil
}

// SetByte has no effect if the address
// is inside the memory.
// Otherwise it returns an error.
func (r *ROM) SetByte(addr uint16, value byte) error {
	if !r.Accepts(addr) {
		return errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return nil
}

// Accepts checks if an address is included in the memory.
func (r *ROM) Accepts(addr uint16) bool {
	return addr < uint16(len(r.rom))
}

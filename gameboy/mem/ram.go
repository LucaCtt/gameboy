package mem

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// RAM represents a generic Random Access Memory.
type RAM struct {
	ram []byte
}

// NewRAM creates a new RAM with addresses from 0x0000
// to the given length.
func NewRAM(len uint16) *RAM {
	return &RAM{make([]byte, len)}
}

// GetByte returns the byte at the given address.
// If the address is outside the memory,
// it will return an error.
func (r *RAM) GetByte(addr uint16) (byte, error) {
	if !r.Accepts(addr) {
		return 0, errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return r.ram[addr], nil
}

// SetByte sets the byte at the given address to the
// given value. Is the address is outside the memory,
// it will return an error.
func (r *RAM) SetByte(addr uint16, value byte) error {
	if !r.Accepts(addr) {
		return errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	r.ram[addr] = value
	return nil
}

// Accepts checks if an address is included in the memory.
func (r *RAM) Accepts(addr uint16) bool {
	return addr < uint16(len(r.ram))
}

package memory

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// RAM represents a generic Random Access Memory.
type RAM struct {
	ram        map[uint16]byte
	start, end uint16
}

// NewRAM creates a new RAM with addresses in
// the given range. The start is inclusive, while
// the end is non-inclusive.
func NewRAM(start, end uint16) *RAM {
	return &RAM{make(map[uint16]byte), start, end}
}

// GetByte returns the byte at the given address.
// If the address is outside the memory space, an
// error will be returned.
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
// given value. Is the address is outside the memory space,
// an error will be returned.
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

// Accepts checks if an address is included in the memory range
func (r *RAM) Accepts(addr uint16) bool {
	return addr >= r.start && addr < r.end
}

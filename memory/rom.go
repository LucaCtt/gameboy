package memory

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// ROM represents a generic Read Only Memory.
type ROM struct {
	rom        map[uint16]byte
	start, end uint16
}

// NewROM creates a new ROM with addresses in
// the given range. The start is inclusive, while
// the end is non-inclusive.
func NewROM(start, end uint16, content map[uint16]byte) *ROM {
	return &ROM{content, start, end}
}

// GetByte returns the byte at the given address.
// If the address is outside the memory space, an
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
// is inside the memory space. Otherwise an
// error will be returned.
func (r *ROM) SetByte(addr uint16, value byte) error {
	if !r.Accepts(addr) {
		return errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return nil
}

// Accepts checks if an address is included in the memory range
func (r *ROM) Accepts(addr uint16) bool {
	return addr >= r.start && addr < r.end
}

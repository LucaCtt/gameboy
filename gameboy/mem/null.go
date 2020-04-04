package mem

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// Null is a memory where writes
// have no effect and reads always return 0x00.
type Null struct {
	len uint16
}

// NewNull creates a new Null with addresses from
// 0x0000 to the given length.
func NewNull(len uint16) *Null {
	return &Null{len}
}

// GetByte returns the byte 0x00 if the address
// is inside the memory. Otherwise it will
// return an error.
func (n *Null) GetByte(addr uint16) (byte, error) {
	if !n.Accepts(addr) {
		return 0, errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return 0, nil
}

// SetByte has no effect if the address
// is inside the memory. Otherwise it
// will return an error.
func (n *Null) SetByte(addr uint16, value byte) error {
	if !n.Accepts(addr) {
		return errors.E(
			fmt.Sprintf("address %v outside of space", addr),
			errors.CodeOutOfRange,
			errors.Memory)
	}
	return nil
}

// Accepts checks if an address is included in the memory.
func (n *Null) Accepts(addr uint16) bool {
	return addr < n.len
}

package memory

import "errors"

// NullMem is a memory space where writes
// have no effect and reads always return 0x00.
type NullMem struct {
	start, end uint16
}

// NewNullMem creates a new NullMemory with addresses in
// the given range. The start is inclusive, while
// the end is non-inclusive.
func NewNullMem(start, end uint16) *NullMem {
	return &NullMem{start, end}
}

// GetByte returns the byte 0x00 if the address
// is inside the memory space. Otherwise an
// error will be returned.
func (n *NullMem) GetByte(addr uint16) (byte, error) {
	if !n.Accepts(addr) {
		return 0, errors.New("")
	}
	return 0, nil
}

// SetByte has no effect if the address
// is inside the memory space. Otherwise an
// error will be returned.
func (n *NullMem) SetByte(addr uint16, value byte) error {
	if !n.Accepts(addr) {
		return errors.New("")
	}
	return nil
}

// Accepts checks if an address is included in the memory space.
func (n *NullMem) Accepts(addr uint16) bool {
	return addr >= n.start && addr < n.end
}

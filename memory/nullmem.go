package memory

// NullMemory is a particular memory space where writes
// have no effect and reads always return 0x00.
type NullMemory struct {
	Start uint16
	End   uint16
}

// GetByte always returns the byte 0x00.
func (n *NullMemory) GetByte(addr uint16) (byte, error) {
	return 0x0000, nil
}

// SetByte has no effect.
func (n *NullMemory) SetByte(addr uint16, value byte) error {
	return nil
}

// Accepts checks if an address is included in the memory space.
func (n *NullMemory) Accepts(addr uint16) bool {
	return addr >= n.Start && addr < n.End
}

// Package cpu implements a complete GameBoy CPU.
package cpu

// Register represents a CPU register.
// As the registers can be used singularly or
// combined to form a 16 bit pseudo-register,
// a common 16 bit representation is used.
type register struct {
	r uint16
}

// Lo returns the lower byte of the register.
func (r *register) Lo() byte {
	return byte(r.r & 0xFF)
}

// Hi returns the higher byte of the register.
func (r *register) Hi() byte {
	return byte(r.r >> 8)
}

// HiLo returns the full value of the register.
func (r *register) HiLo() uint16 {
	return r.r
}

// SetHi sets the value of the higher byte of the register.
func (r *register) SetHi(value byte) {
	r.r = uint16(value)<<8 | (r.r & 0xFF)
}

// SetLo sets the value of the lower byte of the register.
func (r *register) SetLo(value byte) {
	r.r = uint16(value) | (r.r & 0xFF00)
}

// Set sets the full value of the register.
func (r *register) Set(value uint16) {
	r.r = value
}

// CPU represents a GameBoy CPU.
type CPU struct {
	AF, BC, DE, HL, SP, PC register
}

// Init initializes the CPU registers
// to the values set on boot by the GameBoy.
func (c *CPU) Init() {
	c.PC.Set(0x100)
	c.SP.Set(0xFFFE)
}

// setFlag sets the value of the bit in the given position
// in the flag register to the given value.
func (c *CPU) setFlag(p int, b bool) {
	var value uint8
	if b {
		value = 1
	}

	c.AF.SetLo(c.AF.Lo() | value<<p)
}

// Z returns true if the zero flag bit is set.
func (c *CPU) Z() bool {
	return c.AF.Lo()>>7 == 1
}

// N returns true if the subtract flag bit is set.
func (c *CPU) N() bool {
	return c.AF.Lo()>>6 == 1
}

// H returns true if the half carry flag bit is set.
func (c *CPU) H() bool {
	return c.AF.Lo()>>5 == 1
}

// C returns true if the carry flag bit is set.
func (c *CPU) C() bool {
	return c.AF.Lo()>>4 == 1
}

// SetZ sets the value of the zero flag.
func (c *CPU) SetZ(value bool) {
	c.setFlag(7, value)
}

// SetN sets the value of the subtract flag.
func (c *CPU) SetN(value bool) {
	c.setFlag(6, value)
}

// SetH sets the value of the half carry flag.
func (c *CPU) SetH(value bool) {
	c.setFlag(5, value)
}

// SetC sets the value of the carry flag.
func (c *CPU) SetC(value bool) {
	c.setFlag(4, value)
}

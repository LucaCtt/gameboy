// Package cpu implements a complete GameBoy CPU.
package cpu

// CPU represents a GameBoy CPU.
type CPU struct {
	a, b, c, d, f, g, h uint8
	sp, pc              uint16
}

// Init initializes the CPU registers
// to the values set on boot by the GameBoy.
func (c *CPU) Init() {
	c.pc = 0x100
	c.sp = 0xFFFE
}

func setBit(n uint8, p int, b bool) uint8 {
	var value uint8 = 0
	if b {
		value = 1
	}

	n |= value << p
	return n
}

// Z returns true if the zero flag bit is set.
func (c *CPU) Z() bool {
	return c.f>>7 == 1
}

// N returns true if the subtract flag bit is set.
func (c *CPU) N() bool {
	return c.f>>6 == 1
}

// H returns true if the half carry flag bit is set.
func (c *CPU) H() bool {
	return c.f>>5 == 1
}

// C returns true if the carry flag bit is set.
func (c *CPU) C() bool {
	return c.f>>4 == 1
}

// SetZ sets the value of the zero flag.
func (c *CPU) SetZ(z bool) {
	c.f = setBit(c.f, 7, z)
}

// SetN sets the value of the subtract flag.
func (c *CPU) SetN(n bool) {
	c.f = setBit(c.f, 6, n)
}

// SetH sets the value of the half carry flag.
func (c *CPU) SetH(n bool) {
	c.f = setBit(c.f, 5, n)
}

// SetC sets the value of the carry flag.
func (c *CPU) SetC(n bool) {
	c.f = setBit(c.f, 4, n)
}

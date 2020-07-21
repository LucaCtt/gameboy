package cpu

// Registers values set on boot.
const (
	defaultAF uint16 = 0x01B0
	defaultBC uint16 = 0x0013
	defaultDE uint16 = 0x00D8
	defaultHL uint16 = 0x014D
	defaultSP uint16 = 0xFFFE
	defaultPC uint16 = 0x0100
)

// reg represents a CPU register.
// As the registers can be used singularly or
// combined to form a 16 bit pseudo-register,
// a common 16 bit representation is used.
type reg struct {
	// This is the actual value of the register.
	r uint16

	// Used only for the F registry, to prevent
	// updating the lower 4 bits.
	mask uint16
}

// Lo returns the lower byte of the register.
func (r *reg) Lo() byte {
	return byte(r.r & 0xFF)
}

// Hi returns the higher byte of the register.
func (r *reg) Hi() byte {
	return byte(r.r >> 8)
}

// HiLo returns the full value of the register.
func (r *reg) HiLo() uint16 {
	return r.r
}

// SetHi sets the value of the higher byte of the register.
func (r *reg) SetHi(value byte) {
	r.r = uint16(value)<<8 | (r.r & 0xFF)
}

// SetLo sets the value of the lower byte of the register.
func (r *reg) SetLo(value byte) {
	r.r = uint16(value) | (r.r & 0xFF00)
	if r.mask != 0 {
		r.r &= r.mask
	}
}

// Set sets the full value of the register.
func (r *reg) Set(value uint16) {
	r.r = value
	if r.mask != 0 {
		r.r &= r.mask
	}
}

// Regs wraps the Gameboy CPU registers.
type Regs struct {
	AF, BC, DE, HL, SP, PC reg
}

// NewRegs creates a new wrapper that contains the CPU registers.
func NewRegs() *Regs {
	regs := &Regs{AF: reg{mask: 0xFFF0}}

	regs.AF.Set(defaultAF)
	regs.BC.Set(defaultBC)
	regs.DE.Set(defaultDE)
	regs.HL.Set(defaultHL)
	regs.SP.Set(defaultSP)
	regs.PC.Set(defaultPC)

	return regs
}

// Z returns true if the zero flag bit is set.
func (r *Regs) Z() bool {
	return r.AF.Lo()>>7 == 1
}

// N returns true if the subtract flag bit is set.
func (r *Regs) N() bool {
	// Do a bitwise AND of the F register with the 0x40 (0b0100) mask
	// to make sure to ignore the most significant byte.
	return (r.AF.Lo()&0x40)>>6 == 1
}

// H returns true if the half carry flag bit is set.
func (r *Regs) H() bool {
	return (r.AF.Lo()&0x20)>>5 == 1
}

// C returns true if the carry flag bit is set.
func (r *Regs) C() bool {
	return (r.AF.Lo()&0x10)>>4 == 1
}

// SetZ sets the value of the zero flag.
func (r *Regs) SetZ(value bool) {
	r.setFlag(7, value)
}

// SetN sets the value of the subtract flag.
func (r *Regs) SetN(value bool) {
	r.setFlag(6, value)
}

// SetH sets the value of the half carry flag.
func (r *Regs) SetH(value bool) {
	r.setFlag(5, value)
}

// SetC sets the value of the carry flag.
func (r *Regs) SetC(value bool) {
	r.setFlag(4, value)
}

// setFlag sets the value of the bit in the given position
// in the flag register to the given value.
func (r *Regs) setFlag(p int, v bool) {
	if v {
		temp := r.AF.Lo() | (1 << p)
		r.AF.SetLo(temp)
	} else {
		r.AF.SetLo(r.AF.Lo() & ^(1 << p))
	}
}

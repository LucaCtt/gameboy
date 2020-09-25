package cpu

import "github.com/lucactt/gameboy/mem"

type instrUtil struct {
	regs *Regs
	mem  mem.Mem
}

// getByte is a wrapper for getting bytes from mem that panic if an error is returned.
func (u *instrUtil) getByte(addr uint16) byte {
	res, err := u.mem.GetByte(addr)
	if err != nil {
		panic(err)
	}
	return res
}

// getByteAtPC gets the byte in memory at the address obtained
// by adding the given offset to the value of PC.
func (u *instrUtil) getByteAtPC(offset uint16) byte {
	return u.getByte(u.regs.PC.HiLo() + offset)
}

// setByte is a wrapper for setting bytes to mem that panic if an error is returned.
func (u *instrUtil) setByte(addr uint16, value byte) {
	err := u.mem.SetByte(addr, value)
	if err != nil {
		panic(err)
	}
}

// inc8 increments an 8 bit register and also sets the correct flags.
func (u *instrUtil) inc8(original byte, set func(byte)) (int, int) {
	res := original + 1
	set(res)

	u.regs.SetZ(res == 0)
	u.regs.SetN(false)

	// Check if there was a carry from the 3th bit
	// by verifying that the original value of B
	// had 1111 for the least significant 4 bits.
	u.regs.SetH((original & 0x0F) == 0x0F)

	return 1, 4
}

// dec8 decrements an 8 bit register and also sets the correct flags.
func (u *instrUtil) dec8(original byte, set func(byte)) (int, int) {
	res := original - 1
	set(res)

	u.regs.SetZ(res == 0)
	u.regs.SetN(true)
	// Check that there was no carry from 4th bit.
	// If there was set H, otherwise unset it.
	u.regs.SetH((original & 0x0F) == 0x00)

	return 1, 4
}

// inc16 increments a 16 bit register.
func (u *instrUtil) inc16(original uint16, set func(uint16)) (int, int) {
	// Note that 16 bit INC/DEC instructions completely ignore flags,
	// while 8 bit INC/DEC do not.
	set(original + 1)
	return 1, 8
}

// dec16 decrements a 16 bit register.
func (u *instrUtil) dec16(original uint16, set func(uint16)) (int, int) {
	set(original - 1)
	return 1, 8
}

// ld8d8 puts an 8 bit immediate value into an 8 bit register.
func (u *instrUtil) ld8d8(set func(byte)) (int, int) {
	set(u.getByteAtPC(1))
	return 2, 8
}

// add16 puts an 8 bit immediate value into an 8 bit register.
func (u *instrUtil) add16(original, a uint16, set func(uint16)) (int, int) {
	res := int32(original) + int32(a)
	set(uint16(res))

	u.regs.SetN(false)
	u.regs.SetH(int32(original&0xFFF) > (res & 0xFFF))
	u.regs.SetC(res > 0xFFFF)

	return 1, 8
}

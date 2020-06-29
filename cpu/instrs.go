package cpu

import (
	"github.com/lucactt/gameboy/mem"
)

// Instr is a Gameboy CPU instruction, which consists in a function
// that returns the number of clock cycles used.
type Instr func() (int, int)

// InstrSet contains both the non-prefixed and CB-prefixed instructions
// supported by the Gameboy CPU.
type InstrSet struct {
	NoPrefix []Instr
	CBPrefix []Instr
}

// NewInstrSet creates a new instruction set that reads and writes to the
// given registers and memory.
func NewInstrSet(regs *Regs, mem mem.Mem) *InstrSet {
	// Wrappers for getting and setting bytes that panic if an error is returned.
	getByte := func(addr uint16) byte {
		res, err := mem.GetByte(addr)
		if err != nil {
			panic(err)
		}
		return res
	}
	setByte := func(addr uint16, value byte) {
		err := mem.SetByte(addr, value)
		if err != nil {
			panic(err)
		}
	}

	// Gets the byte in memory at the address obtained
	// by adding the given offset to the value of PC.
	getByteAtPC := func(offset uint16) byte {
		return getByte(regs.PC.HiLo() + offset)
	}

	return &InstrSet{
		NoPrefix: []Instr{
			func() (int, int) {
				// 0x00 - NOP
				return 1, 4
			},
			func() (int, int) {
				// 0x01 - LD BC,d16

				// The 16 bit value to be loaded in BC is found
				// in the two bytes at mem[PC+1] and mem[PC+2],
				// where the former is the least significant byte.
				regs.BC.SetLo(getByteAtPC(1))
				regs.BC.SetHi(getByteAtPC(2))
				return 3, 12
			},
			func() (int, int) {
				// 0x02 - LD (BC),A
				setByte(regs.BC.HiLo(), regs.AF.Hi())
				return 1, 8
			},
			func() (int, int) {
				// 0x03 - INC BC

				// Note that 16 bit INC/DEC instructions completely ignore flags,
				// while 8 bit INC/DEC do not.
				regs.BC.Set(regs.BC.HiLo() + 1)
				return 1, 8
			},
			func() (int, int) {
				// 0x04 - INC B
				original := regs.BC.Hi()
				regs.BC.SetHi(original + 1)

				regs.SetZ(regs.BC.Hi() == 0)
				regs.SetN(false)

				// Check if there was a carry from the 3th bit
				// by verifying that the original value of B
				// had 1111 for the least significant 4 bits.
				regs.SetH((original & 0x0F) == 0x0F)

				return 1, 4
			},
			func() (int, int) {
				// 0x05 - DEC B
				original := regs.BC.Hi()
				regs.BC.SetHi(original + 1)

				regs.SetZ(regs.BC.Hi() == 0)
				regs.SetN(true)
				regs.SetH((original & 0x0F) == 0x0F)

				return 1, 4
			},
			func() (int, int) {
				// 0x06 - LD B,d8
				regs.BC.SetHi(getByteAtPC(1))
				return 2, 8
			},
			func() (int, int) {
				// 0x07 - RLCA
				original := regs.AF.Hi()
				regs.AF.SetHi(original << 1)

				regs.SetZ(regs.AF.Hi() == 0)
				regs.SetN(false)
				regs.SetH(false)

				// Put the old value of the 7th bit of A in flag C.
				regs.SetC((original & (1 << 7)) != 0)
				return 1, 4
			},
			func() (int, int) {
				// 0x08 - LD (a16),SP

				// Compose the 16 bit address from the two 8 bit parts.
				addr := uint16(getByteAtPC(2))<<8 | uint16(getByteAtPC(1))
				setByte(addr, regs.SP.Hi())
				setByte(addr+1, regs.SP.Lo())
				return 3, 20
			},
		},
	}
}

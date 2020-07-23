package cpu

import (
	"github.com/lucactt/gameboy/mem"
)

// Instr is a Gameboy CPU instruction, which consists in a function
// that returns the number of bytes read and number of clock cycles used.
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
	util := &instrUtil{regs, mem}

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
				regs.BC.SetLo(util.getByteAtPC(1))
				regs.BC.SetHi(util.getByteAtPC(2))
				return 3, 12
			},
			func() (int, int) {
				// 0x02 - LD (BC),A
				util.setByte(regs.BC.HiLo(), regs.AF.Hi())
				return 1, 8
			},
			func() (int, int) {
				// 0x03 - INC BC
				return util.inc16(regs.BC.HiLo(), func(res uint16) { regs.BC.Set(res) })
			},
			func() (int, int) {
				// 0x04 - INC B
				return util.inc8(regs.BC.Hi(), func(res byte) { regs.BC.SetHi(res) })
			},
			func() (int, int) {
				// 0x05 - DEC B
				return util.dec8(regs.BC.Hi(), func(res byte) { regs.BC.SetHi(res) })
			},
			func() (int, int) {
				// 0x06 - LD B,d8
				regs.BC.SetHi(util.getByteAtPC(1))
				return 2, 8
			},
			func() (int, int) {
				// 0x07 - RLCA
				original := regs.AF.Hi()
				regs.AF.SetHi((original << 1) | (original >> 7))

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
				addr := uint16(util.getByteAtPC(2))<<8 | uint16(util.getByteAtPC(1))
				util.setByte(addr, regs.SP.Lo())
				util.setByte(addr+1, regs.SP.Hi())
				return 3, 20
			},
			func() (int, int) {
				// 0x09 - ADD HL,BC

				original := regs.HL.HiLo()
				res := int32(original) + int32(regs.BC.HiLo())
				regs.HL.Set(uint16(res))

				regs.SetN(false)
				regs.SetH(int32(original&0xFFF) > (res & 0xFFF))
				regs.SetC(res > 0xFFFF)

				return 1, 8
			},
			func() (int, int) {
				// 0x0A - LD A,(BC)
				regs.AF.SetHi(util.getByte(regs.BC.HiLo()))
				return 1, 8
			},
			func() (int, int) {
				// 0x0B - DEC BC
				return util.dec16(regs.BC.HiLo(), func(res uint16) { regs.BC.Set(res) })
			},
			func() (int, int) {
				// 0x0C - INC C
				return util.inc8(regs.BC.Lo(), func(res byte) { regs.BC.SetLo(res) })
			},
			func() (int, int) {
				// 0x0D - DEC C
				return util.dec8(regs.BC.Lo(), func(res byte) { regs.BC.SetLo(res) })
			},
			func() (int, int) {
				// 0x0E - LD C,d8
				regs.BC.SetLo(util.getByteAtPC(1))
				return 2, 8
			},
			func() (int, int) {
				// 0x0F - RRCA
				original := regs.AF.Hi()
				regs.AF.SetHi((original >> 1) | (original << 7))

				regs.SetZ(regs.AF.Hi() == 0)
				regs.SetN(false)
				regs.SetH(false)

				// Put the old value of the 0th bit of A in flag C.
				regs.SetC((original & 0x01) != 0)
				return 1, 4
			},
		},
	}
}

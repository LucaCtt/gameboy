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
func NewInstrSet(regs *Regs, mem mem.Mem, stateMgr *StateMgr) *InstrSet {
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
				return util.inc16(regs.BC.HiLo(), regs.BC.Set)
			},
			func() (int, int) {
				// 0x04 - INC B
				return util.inc8(regs.BC.Hi(), regs.BC.SetHi)
			},
			func() (int, int) {
				// 0x05 - DEC B
				return util.dec8(regs.BC.Hi(), regs.BC.SetHi)
			},
			func() (int, int) {
				// 0x06 - LD B,d8
				return util.ld8d8(regs.BC.SetHi)
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
				return util.add16(regs.HL.HiLo(), regs.BC.HiLo(), func(res uint16) { regs.HL.Set(res) })
			},
			func() (int, int) {
				// 0x0A - LD A,(BC)
				regs.AF.SetHi(util.getByte(regs.BC.HiLo()))
				return 1, 8
			},
			func() (int, int) {
				// 0x0B - DEC BC
				return util.dec16(regs.BC.HiLo(), regs.BC.Set)
			},
			func() (int, int) {
				// 0x0C - INC C
				return util.inc8(regs.BC.Lo(), regs.BC.SetLo)
			},
			func() (int, int) {
				// 0x0D - DEC C
				return util.dec8(regs.BC.Lo(), regs.BC.SetLo)
			},
			func() (int, int) {
				// 0x0E - LD C,d8
				return util.ld8d8(regs.BC.SetLo)
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
			func() (int, int) {
				// 0x10 - STOP
				stateMgr.SetState(Stopped)
				return 2, 4
			},
			func() (int, int) {
				// 0x11 - LD DE,d16
				regs.DE.SetLo(util.getByteAtPC(1))
				regs.DE.SetHi(util.getByteAtPC(2))
				return 3, 12
			},
			func() (int, int) {
				// 0x12 - LD (DE),A
				util.setByte(regs.DE.HiLo(), regs.AF.Hi())
				return 1, 8
			},
			func() (int, int) {
				// 0x13 - INC DE
				return util.inc16(regs.DE.HiLo(), regs.DE.Set)
			},
			func() (int, int) {
				// 0x14 - INC D
				return util.inc8(regs.DE.Hi(), regs.DE.SetHi)
			},
			func() (int, int) {
				// 0x15 - DEC D
				return util.dec8(regs.DE.Hi(), regs.DE.SetHi)
			},
			func() (int, int) {
				// 0x16 - LD D,d8
				return util.ld8d8(regs.DE.SetHi)
			},
			func() (int, int) {
				// 0x17 - RLA
				original := regs.AF.Hi()
				var carry byte
				if regs.C() {
					carry = 1
				}
				regs.AF.SetHi((original << 1) + carry)

				regs.SetZ(regs.AF.Hi() == 0)
				regs.SetN(false)
				regs.SetH(false)

				// Put the old value of the 7th bit of A in flag C.
				regs.SetC((original & (1 << 7)) != 0)
				return 1, 4
			},
			func() (int, int) {
				// 0x18 - JR r8
				regs.PC.Set(regs.PC.HiLo() + uint16(util.getByteAtPC(1)))
				return 2, 12
			},
			func() (int, int) {
				// 0x19 - ADD HL,DE
				return util.add16(regs.HL.HiLo(), regs.DE.HiLo(), func(res uint16) { regs.HL.Set(res) })
			},
			func() (int, int) {
				// 0x1A - LD A,(DE)
				regs.AF.SetHi(util.getByte(regs.DE.HiLo()))
				return 1, 8
			},
			func() (int, int) {
				// 0x1B - DEC DE
				return util.dec16(regs.DE.HiLo(), regs.DE.Set)
			},
			func() (int, int) {
				// 0x1C - INC E
				return util.inc8(regs.DE.Lo(), regs.DE.SetLo)
			},
			func() (int, int) {
				// 0x1D - DEC E
				return util.dec8(regs.DE.Lo(), regs.DE.SetLo)
			},
			func() (int, int) {
				// 0x1E - LD E,d8
				return util.ld8d8(regs.DE.SetLo)
			},
			func() (int, int) {
				// 0x1F - RRA
				original := regs.AF.Hi()
				var carry byte
				if regs.C() {
					carry = 1
				}
				regs.AF.SetHi((original >> 1) | (carry << 7))

				regs.SetZ(regs.AF.Hi() == 0)
				regs.SetN(false)
				regs.SetH(false)

				// Put the old value of the 0th bit of A in flag C.
				regs.SetC((original & 0x01) != 0)
				return 1, 4
			},
			func() (int, int) {
				// 0x20 - JR NZ,r8
				if !regs.Z() {
					regs.PC.Set(regs.PC.HiLo() + uint16(util.getByteAtPC(1)))
				}
				return 2, 8
			},
			func() (int, int) {
				// 0x21 - LD HL,d16

				regs.HL.SetLo(util.getByteAtPC(1))
				regs.HL.SetHi(util.getByteAtPC(2))
				return 3, 12
			},
			func() (int, int) {
				// 0x22 - LD (HL+),A
				util.setByte(regs.HL.HiLo(), regs.AF.Hi())
				util.inc16(regs.HL.HiLo(), regs.HL.Set)
				return 1, 8
			},
			func() (int, int) {
				// 0x23 - INC HL
				return util.inc16(regs.HL.HiLo(), regs.HL.Set)
			},
			func() (int, int) {
				// 0x24 - INC H
				return util.inc8(regs.HL.Hi(), regs.HL.SetHi)
			},
			func() (int, int) {
				// 0x25 - DEC H
				return util.dec8(regs.HL.Hi(), regs.HL.SetHi)
			},
			func() (int, int) {
				// 0x26 - LD H,d8
				return util.ld8d8(regs.HL.SetHi)
			},
			func() (int, int) {
				// 0x27 - DAA

				// Copied and pasted from goboy cause fuck this shit
				if !regs.N() {
					if regs.C() || regs.AF.Hi() > 0x99 {
						regs.AF.SetHi(regs.AF.Hi() + 0x60)
						regs.SetC(true)
					}
					if regs.H() || regs.AF.Hi()&0xF > 0x9 {
						regs.AF.SetHi(regs.AF.Hi() + 0x06)
						regs.SetH(false)
					}
				} else if regs.C() && regs.H() {
					regs.AF.SetHi(regs.AF.Hi() + 0x9A)
					regs.SetH(false)
				} else if regs.C() {
					regs.AF.SetHi(regs.AF.Hi() + 0xA0)
				} else if regs.H() {
					regs.AF.SetHi(regs.AF.Hi() + 0xFA)
					regs.SetH(false)
				}
				regs.SetZ(regs.AF.Hi() == 0)

				return 1, 4
			},
			func() (int, int) {
				// 0x28 - JR Z, r8
				if regs.Z() {
					regs.PC.Set(regs.PC.HiLo() + uint16(util.getByteAtPC(1)))
				}
				return 2, 8
			},
			func() (int, int) {
				// 0x29 - ADD HL,HL
				return util.add16(regs.HL.HiLo(), regs.HL.HiLo(), func(res uint16) { regs.HL.Set(res) })
			},
			func() (int, int) {
				// 0x2A - LD A,(HL+)
				regs.AF.SetHi(util.getByte(regs.HL.HiLo()))
				util.inc16(regs.HL.HiLo(), regs.HL.Set)
				return 1, 8
			},
			func() (int, int) {
				// 0x2B - DEC HL
				return util.dec16(regs.HL.HiLo(), regs.HL.Set)
			},
			func() (int, int) {
				// 0x2C - INC L
				return util.inc8(regs.HL.Lo(), regs.HL.SetLo)
			},
			func() (int, int) {
				// 0x2D - DEC L
				return util.dec8(regs.HL.Lo(), regs.HL.SetLo)
			},
			func() (int, int) {
				// 0x2E - LD L,d8
				return util.ld8d8(regs.HL.SetLo)
			},
			func() (int, int) {
				// 0x2F - CPL
				regs.AF.SetHi(^regs.AF.Hi())
				regs.SetN(true)
				regs.SetH(true)

				return 1, 4
			},
		},
	}
}

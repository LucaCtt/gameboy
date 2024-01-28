package cpu

import (
	"testing"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func Test_NewInstrSet(t *testing.T) {
	t.Run("no prefix", func(t *testing.T) {
		t.Run("NOP", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(0)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			len, cycles := set.NoPrefix[0x00]()

			assert.Equal(t, regs, NewRegs())
			assert.Equal(t, ram, mem.NewRAM(0))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 4)
		})

		t.Run("LD BC,d16", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 3)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x01)
			ram.SetByte(regs.PC.HiLo()+2, 0x11)

			len, cycles := set.NoPrefix[0x01]()

			assert.Equal(t, regs.BC.HiLo(), uint16(0x1101))
			assert.Equal(t, len, 3)
			assert.Equal(t, cycles, 12)
		})

		t.Run("LD (BC),A", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.BC.Set(0x0001)
			ram.SetByte(regs.PC.HiLo()+1, 0x01)

			len, cycles := set.NoPrefix[0x02]()

			got, _ := ram.GetByte(regs.BC.HiLo())
			assert.Equal(t, got, byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("INC BC", func(t *testing.T) {
			testInc16(t, 0x03,
				func(regs *Regs) uint16 { return regs.BC.HiLo() },
				func(regs *Regs, v uint16) { regs.BC.Set(v) })
		})

		t.Run("INC B", func(t *testing.T) {
			testInc8(t, 0x04,
				func(regs *Regs) byte { return regs.BC.Hi() },
				func(regs *Regs, v byte) { regs.BC.SetHi(v) })
		})

		t.Run("DEC B", func(t *testing.T) {
			testDec8(t, 0x05,
				func(regs *Regs) byte { return regs.BC.Hi() },
				func(regs *Regs, v byte) { regs.BC.SetHi(v) })
		})

		t.Run("LD B,d8", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x10)

			len, cycles := set.NoPrefix[0x06]()

			assert.Equal(t, regs.BC.Hi(), byte(0x10))
			assert.Equal(t, len, 2)
			assert.Equal(t, cycles, 8)
		})

		t.Run("RLCA", func(t *testing.T) {
			t.Run("msb is 0", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x01)

				len, cycles := set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0x02))
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
			})

			t.Run("msb is 1", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0xF0) //0b11110000

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0xE1)) // should be 0b11100001
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

			t.Run("result is zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x00)

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
			})
		})

		t.Run("LD (d16),SP", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(258)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.SP.Set(0x0110)
			regs.PC.Set(0x0000)
			ram.SetByte(regs.PC.HiLo()+1, 0x00)
			ram.SetByte(regs.PC.HiLo()+2, 0x01)

			len, cycles := set.NoPrefix[0x08]()

			lo, _ := ram.GetByte(0x0100)
			hi, _ := ram.GetByte(0x0101)

			assert.Equal(t, lo, byte(0x10))
			assert.Equal(t, hi, byte(0x01))
			assert.Equal(t, len, 3)
			assert.Equal(t, cycles, 20)
		})

		t.Run("ADD HL,BC", func(t *testing.T) {
			testAdd16(t, 0x09,
				func(regs *Regs) uint16 { return regs.HL.HiLo() },
				func(regs *Regs, s, a uint16) {
					regs.HL.Set(s)
					regs.BC.Set(a)
				})
		})

		t.Run("LD A,(BC)", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(1)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.BC.Set(0x0000)
			ram.SetByte(regs.BC.HiLo(), 0x01)

			len, cycles := set.NoPrefix[0x0A]()

			assert.Equal(t, regs.AF.Hi(), byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("DEC BC", func(t *testing.T) {
			testDec16(t, 0x0B,
				func(regs *Regs) uint16 { return regs.BC.HiLo() },
				func(regs *Regs, v uint16) { regs.BC.Set(v) })
		})

		t.Run("INC C", func(t *testing.T) {
			testInc8(t, 0x0C,
				func(regs *Regs) byte { return regs.BC.Lo() },
				func(regs *Regs, v byte) { regs.BC.SetLo(v) })
		})

		t.Run("DEC C", func(t *testing.T) {
			testDec8(t, 0x0D,
				func(regs *Regs) byte { return regs.BC.Lo() },
				func(regs *Regs, v byte) { regs.BC.SetLo(v) })
		})

		t.Run("LD C,d8", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x10)

			len, cycles := set.NoPrefix[0x0E]()

			assert.Equal(t, regs.BC.Lo(), byte(0x10))
			assert.Equal(t, len, 2)
			assert.Equal(t, cycles, 8)
		})

		t.Run("RRCA", func(t *testing.T) {
			t.Run("msb is 0", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x01)

				len, cycles := set.NoPrefix[0x0F]()

				assert.Equal(t, regs.AF.Hi(), byte(0x80))
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

			t.Run("msb is 1", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x0F) //0b00001111

				set.NoPrefix[0x0F]()

				assert.Equal(t, regs.AF.Hi(), byte(0x87)) // should be 0b10000111
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

			t.Run("result is zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x00)

				set.NoPrefix[0x0F]()

				assert.Equal(t, regs.AF.Hi(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
			})
		})

		t.Run("STOP", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(0)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			len, cycles := set.NoPrefix[0x10]()

			assert.Equal(t, stateMgr.current, Stopped)
			assert.Equal(t, len, 2)
			assert.Equal(t, cycles, 4)
		})

		t.Run("LD DE,d16", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 3)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x01)
			ram.SetByte(regs.PC.HiLo()+2, 0x11)

			len, cycles := set.NoPrefix[0x11]()

			assert.Equal(t, regs.DE.HiLo(), uint16(0x1101))
			assert.Equal(t, len, 3)
			assert.Equal(t, cycles, 12)
		})

		t.Run("LD (DE),A", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.DE.Set(0x0001)
			ram.SetByte(regs.PC.HiLo()+1, 0x01)

			len, cycles := set.NoPrefix[0x12]()

			got, _ := ram.GetByte(regs.DE.HiLo())
			assert.Equal(t, got, byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("INC DE", func(t *testing.T) {
			testInc16(t, 0x13,
				func(regs *Regs) uint16 { return regs.DE.HiLo() },
				func(regs *Regs, v uint16) { regs.DE.Set(v) })
		})

		t.Run("INC D", func(t *testing.T) {
			testInc8(t, 0x14,
				func(regs *Regs) byte { return regs.DE.Hi() },
				func(regs *Regs, v byte) { regs.DE.SetHi(v) })
		})

		t.Run("DEC D", func(t *testing.T) {
			testDec8(t, 0x15,
				func(regs *Regs) byte { return regs.DE.Hi() },
				func(regs *Regs, v byte) { regs.DE.SetHi(v) })
		})

		t.Run("LD D,d8", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x10)

			len, cycles := set.NoPrefix[0x16]()

			assert.Equal(t, regs.DE.Hi(), byte(0x10))
			assert.Equal(t, len, 2)
			assert.Equal(t, cycles, 8)
		})

		t.Run("RLA", func(t *testing.T) {
			//TODO
			t.Run("msb is 0", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x01)

				len, cycles := set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0x02))
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
			})

			t.Run("msb is 1", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0xF0) //0b11110000

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0xE1)) // should be 0b11100001
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

			t.Run("result is zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				stateMgr := NewStateMgr()
				set := NewInstrSet(regs, ram, stateMgr)

				regs.AF.SetHi(0x00)

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
			})
		})

		t.Run("JR r8", func(t *testing.T) {
			//TODO
			regs := NewRegs()
			ram := mem.NewRAM(258)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.SP.Set(0x0110)
			regs.PC.Set(0x0000)
			ram.SetByte(regs.PC.HiLo()+1, 0x00)
			ram.SetByte(regs.PC.HiLo()+2, 0x01)

			len, cycles := set.NoPrefix[0x08]()

			lo, _ := ram.GetByte(0x0100)
			hi, _ := ram.GetByte(0x0101)

			assert.Equal(t, lo, byte(0x10))
			assert.Equal(t, hi, byte(0x01))
			assert.Equal(t, len, 3)
			assert.Equal(t, cycles, 20)
		})

		t.Run("ADD HL,DE", func(t *testing.T) {
			testAdd16(t, 0x19,
				func(regs *Regs) uint16 { return regs.HL.HiLo() },
				func(regs *Regs, s, a uint16) {
					regs.HL.Set(s)
					regs.DE.Set(a)
				})
		})

		t.Run("LD A,(DE)", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(1)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			regs.DE.Set(0x0000)
			ram.SetByte(regs.DE.HiLo(), 0x01)

			len, cycles := set.NoPrefix[0x1A]()

			assert.Equal(t, regs.AF.Hi(), byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("DEC DE", func(t *testing.T) {
			testDec16(t, 0x1B,
				func(regs *Regs) uint16 { return regs.DE.HiLo() },
				func(regs *Regs, v uint16) { regs.DE.Set(v) })
		})

		t.Run("INC E", func(t *testing.T) {
			testInc8(t, 0x1C,
				func(regs *Regs) byte { return regs.DE.Lo() },
				func(regs *Regs, v byte) { regs.DE.SetLo(v) })
		})

		t.Run("DEC E", func(t *testing.T) {
			testDec8(t, 0x1D,
				func(regs *Regs) byte { return regs.DE.Lo() },
				func(regs *Regs, v byte) { regs.DE.SetLo(v) })
		})

		t.Run("LD E,d8", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 2)
			stateMgr := NewStateMgr()
			set := NewInstrSet(regs, ram, stateMgr)

			ram.SetByte(regs.PC.HiLo()+1, 0x10)

			len, cycles := set.NoPrefix[0x1E]()

			assert.Equal(t, regs.DE.Lo(), byte(0x10))
			assert.Equal(t, len, 2)
			assert.Equal(t, cycles, 8)
		})
	})
}

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
			set := NewInstrSet(regs, ram)

			len, cycles := set.NoPrefix[0x00]()

			assert.Equal(t, regs, NewRegs())
			assert.Equal(t, ram, mem.NewRAM(0))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 4)
		})

		t.Run("LD BC,d16", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 3)
			set := NewInstrSet(regs, ram)

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
			set := NewInstrSet(regs, ram)

			regs.BC.Set(0x0001)
			ram.SetByte(regs.PC.HiLo()+1, 0x01)

			len, cycles := set.NoPrefix[0x02]()

			got, _ := ram.GetByte(regs.BC.HiLo())
			assert.Equal(t, got, byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("INC BC", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(0)
			set := NewInstrSet(regs, ram)

			regs.BC.Set(0x0001)

			len, cycles := set.NoPrefix[0x03]()

			assert.Equal(t, regs.BC.HiLo(), uint16(0x0002))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("INC B", func(t *testing.T) {
			t.Run("no carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0x01)

				len, cycles := set.NoPrefix[0x04]()

				assert.Equal(t, regs.BC.Hi(), byte(0x02))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
			})

			t.Run("carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0x0F)

				set.NoPrefix[0x04]()

				assert.Equal(t, regs.BC.Hi(), byte(0x10))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
			})

			t.Run("overflow", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0xFF)

				set.NoPrefix[0x04]()

				assert.Equal(t, regs.BC.Hi(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
			})
		})

		t.Run("DEC B", func(t *testing.T) {
			t.Run("no carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0x02)

				len, cycles := set.NoPrefix[0x05]()

				assert.Equal(t, regs.BC.Hi(), byte(0x01))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
			})

			t.Run("carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0x10)

				set.NoPrefix[0x05]()

				assert.Equal(t, regs.BC.Hi(), byte(0x0F))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), true)
			})

			t.Run("zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetHi(0x01)

				set.NoPrefix[0x05]()

				assert.Equal(t, regs.BC.Hi(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), false)
			})
		})

		t.Run("LD B,d8", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(regs.PC.HiLo() + 2)
			set := NewInstrSet(regs, ram)

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
				set := NewInstrSet(regs, ram)

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
				set := NewInstrSet(regs, ram)

				regs.AF.SetHi(0xF0) //0b11110000

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0xE0)) // should be 0b11100000
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

			t.Run("result is zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.AF.SetHi(0x80) //0b10000000

				set.NoPrefix[0x07]()

				assert.Equal(t, regs.AF.Hi(), byte(0x00)) // should be 0b11100000
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), true)
			})

		})

		t.Run("LD (d16),SP", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(258)
			set := NewInstrSet(regs, ram)

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
			t.Run("no carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.HL.Set(0x0001)
				regs.BC.Set(0x0001)

				len, cycles := set.NoPrefix[0x09]()

				assert.Equal(t, regs.HL.HiLo(), uint16(0x0002))
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, regs.C(), false)
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 8)
			})

			t.Run("half carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.HL.Set(0x0FFF)
				regs.BC.Set(0x0001)

				set.NoPrefix[0x09]()

				assert.Equal(t, regs.HL.HiLo(), uint16(0x1000))
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
				assert.Equal(t, regs.C(), false)
			})

			t.Run("carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.HL.Set(0xFFFF)
				regs.BC.Set(0x0001)

				set.NoPrefix[0x09]()

				assert.Equal(t, regs.HL.HiLo(), uint16(0x0000))
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
				assert.Equal(t, regs.C(), true)
			})
		})

		t.Run("LD A,(BC)", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(1)
			set := NewInstrSet(regs, ram)

			regs.BC.Set(0x0000)
			ram.SetByte(regs.BC.HiLo(), 0x01)

			len, cycles := set.NoPrefix[0x0A]()

			assert.Equal(t, regs.AF.Hi(), byte(0x01))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("DEC BC", func(t *testing.T) {
			regs := NewRegs()
			ram := mem.NewRAM(0)
			set := NewInstrSet(regs, ram)

			regs.BC.Set(0x0001)

			len, cycles := set.NoPrefix[0x0B]()

			assert.Equal(t, regs.BC.HiLo(), uint16(0x0000))
			assert.Equal(t, len, 1)
			assert.Equal(t, cycles, 8)
		})

		t.Run("INC C", func(t *testing.T) {
			t.Run("no carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0x01)

				len, cycles := set.NoPrefix[0x0C]()

				assert.Equal(t, regs.BC.Lo(), byte(0x02))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
			})

			t.Run("carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0x0F)

				set.NoPrefix[0x0C]()

				assert.Equal(t, regs.BC.Lo(), byte(0x10))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
			})

			t.Run("overflow", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0xFF)

				set.NoPrefix[0x0C]()

				assert.Equal(t, regs.BC.Lo(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), false)
				assert.Equal(t, regs.H(), true)
			})
		})

		t.Run("DEC C", func(t *testing.T) {
			t.Run("no carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0x02)

				len, cycles := set.NoPrefix[0x0D]()

				assert.Equal(t, regs.BC.Lo(), byte(0x01))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), false)
				assert.Equal(t, len, 1)
				assert.Equal(t, cycles, 4)
			})

			t.Run("carry", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0x10)

				set.NoPrefix[0x0D]()

				assert.Equal(t, regs.BC.Lo(), byte(0x0F))
				assert.Equal(t, regs.Z(), false)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), true)
			})

			t.Run("zero", func(t *testing.T) {
				regs := NewRegs()
				ram := mem.NewRAM(0)
				set := NewInstrSet(regs, ram)

				regs.BC.SetLo(0x01)

				set.NoPrefix[0x0D]()

				assert.Equal(t, regs.BC.Lo(), byte(0x00))
				assert.Equal(t, regs.Z(), true)
				assert.Equal(t, regs.N(), true)
				assert.Equal(t, regs.H(), false)
			})
		})
	})
}

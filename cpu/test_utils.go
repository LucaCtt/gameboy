package cpu

import (
	"testing"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func testInc8(t *testing.T, opcode byte, getR func(*Regs) byte, setR func(*Regs, byte)) {
	t.Helper()

	t.Run("no carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x01)

		len, cycles := set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x02))
		assert.Equal(t, regs.Z(), false)
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), false)
		assert.Equal(t, len, 1)
		assert.Equal(t, cycles, 4)
	})

	t.Run("carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x0F)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x10))
		assert.Equal(t, regs.Z(), false)
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), true)
	})

	t.Run("overflow", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0xFF)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x00))
		assert.Equal(t, regs.Z(), true)
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), true)
	})
}

func testDec8(t *testing.T, opcode byte, getR func(*Regs) byte, setR func(*Regs, byte)) {
	t.Helper()

	t.Run("no carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x02)

		len, cycles := set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x01))
		assert.Equal(t, regs.Z(), false)
		assert.Equal(t, regs.N(), true)
		assert.Equal(t, regs.H(), false)
		assert.Equal(t, len, 1)
		assert.Equal(t, cycles, 4)
	})

	t.Run("carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x10)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x0F))
		assert.Equal(t, regs.Z(), false)
		assert.Equal(t, regs.N(), true)
		assert.Equal(t, regs.H(), true)
	})

	t.Run("zero", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x01)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), byte(0x00))
		assert.Equal(t, regs.Z(), true)
		assert.Equal(t, regs.N(), true)
		assert.Equal(t, regs.H(), false)
	})
}

func testInc16(t *testing.T, opcode byte, getR func(*Regs) uint16, setR func(*Regs, uint16)) {
	t.Helper()

	regs := NewRegs()
	ram := mem.NewRAM(0)
	stateMgr := NewStateMgr()
	set := NewInstrSet(regs, ram, stateMgr)

	setR(regs, 0x0001)

	len, cycles := set.NoPrefix[opcode]()

	assert.Equal(t, getR(regs), uint16(0x0002))
	assert.Equal(t, len, 1)
	assert.Equal(t, cycles, 8)
}

func testDec16(t *testing.T, opcode byte, getR func(*Regs) uint16, setR func(*Regs, uint16)) {
	t.Helper()

	regs := NewRegs()
	ram := mem.NewRAM(0)
	stateMgr := NewStateMgr()
	set := NewInstrSet(regs, ram, stateMgr)

	setR(regs, 0x0001)

	len, cycles := set.NoPrefix[opcode]()

	assert.Equal(t, getR(regs), uint16(0x0000))
	assert.Equal(t, len, 1)
	assert.Equal(t, cycles, 8)
}

func testAdd16(t *testing.T, opcode byte, getR func(*Regs) uint16, setR func(*Regs, uint16, uint16)) {
	t.Helper()

	t.Run("no carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x0001, 0x0001)

		len, cycles := set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), uint16(0x0002))
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), false)
		assert.Equal(t, regs.C(), false)
		assert.Equal(t, len, 1)
		assert.Equal(t, cycles, 8)
	})

	t.Run("half carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0x0FFF, 0x0001)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), uint16(0x1000))
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), true)
		assert.Equal(t, regs.C(), false)
	})

	t.Run("carry", func(t *testing.T) {
		regs := NewRegs()
		ram := mem.NewRAM(0)
		stateMgr := NewStateMgr()
		set := NewInstrSet(regs, ram, stateMgr)

		setR(regs, 0xFFFF, 0x0001)

		set.NoPrefix[opcode]()

		assert.Equal(t, getR(regs), uint16(0x0000))
		assert.Equal(t, regs.N(), false)
		assert.Equal(t, regs.H(), true)
		assert.Equal(t, regs.C(), true)
	})
}

func testld8(t *testing.T, opcode byte, getR func(*Regs) byte, setR func(*Regs, byte)) {
	t.Helper()

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
	regs := NewRegs()
	ram := mem.NewRAM(0)
	stateMgr := NewStateMgr()
	set := NewInstrSet(regs, ram, stateMgr)

	setR(regs, 0x0001)

	len, cycles := set.NoPrefix[opcode]()

	assert.Equal(t, getR(regs), uint16(0x0000))
	assert.Equal(t, len, 1)
	assert.Equal(t, cycles, 8)
}

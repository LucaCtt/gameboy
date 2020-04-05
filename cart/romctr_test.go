package cart

import (
	"testing"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func Test_newROMCtr(t *testing.T) {
	t.Run("ROM size is too small", func(t *testing.T) {
		bytes := make([]byte, 0)
		rom := mem.NewROM(bytes)

		_, err := NewROMCtr(rom)
		assert.Err(t, err, true)
	})

	t.Run("ROM size is too big", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		_, err := NewROMCtr(rom)
		assert.Err(t, err, true)
	})

	t.Run("ROM size is big enough", func(t *testing.T) {
		// The length of this slice must be at least romEnd + 1, and not just romEnd
		// because that would create a slice with addresses between 0 and romEnd - 1.
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)

		_, err := NewROMCtr(rom)
		assert.Err(t, err, false)
	})
}

func Test_ROMCtr_GetByte(t *testing.T) {
	t.Run("RAM addr", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		// Set the ram banks flag byte to a value that indicates at least 1 bank.
		bytes[ramSize] = ramBank1

		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		// The ram is created by the controller, so it can be accessed
		// only by using SetByte.
		err := ctr.SetByte(ramStart, 0x10)
		assert.Err(t, err, false)

		got, err := ctr.GetByte(ramStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x10))
	})

	t.Run("ROM addr", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		// Here the value of the rom byte is set directly.
		bytes[romStart] = 0x11

		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		got, err := ctr.GetByte(romStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("Invalid addr", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)

		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		_, err := ctr.GetByte(romEnd + 1)
		assert.Err(t, err, true)
	})
}

func Test_ROMCtr_SetByte(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		bytes[ramSize] = ramBank1

		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		err := ctr.SetByte(ramEnd, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(ramEnd)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("ROM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		err := ctr.SetByte(romEnd, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(romEnd)
		assert.Equal(t, got, byte(0x00))
	})

	t.Run("Invalid address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		err := ctr.SetByte(romEnd+1, 0x11)
		assert.Err(t, err, true)
	})
}

func Test_ROMCtr_Accepts(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		bytes[ramSize] = ramBank1

		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		got := ctr.Accepts(ramStart)
		assert.Equal(t, got, true)

		got = ctr.Accepts(ramEnd)
		assert.Equal(t, got, true)
	})

	t.Run("ROM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)
		ctr, _ := NewROMCtr(rom)

		got := ctr.Accepts(romStart)
		assert.Equal(t, got, true)

		got = ctr.Accepts(romEnd)
		assert.Equal(t, got, true)
	})
}

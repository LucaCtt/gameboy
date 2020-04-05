package cart

import (
	"testing"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func Test_newROMCtr(t *testing.T) {
	t.Run("rom size is too small", func(t *testing.T) {
		bytes := make([]byte, 0)
		rom := mem.NewROM(bytes)

		got, err := NewROM(rom)
		assert.Err(t, err, true)
		assert.NotEqual(t, got, nil)
	})

	t.Run("rom size is big enough", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		got, err := NewROM(rom)
		assert.Err(t, err, false)
		assert.NotEqual(t, got, nil)
	})
}

func Test_romCtr_GetByte(t *testing.T) {
	t.Run("RAM addr", func(t *testing.T) {
		// The length of this slice must be at least romEnd + 1, and not just romEnd
		// because that would create a slice with addresses between 0 and romEnd - 1.
		bytes := make([]byte, romEnd+1)
		// Set the ram banks flag byte to a value that indicates at least 1 bank.
		bytes[ramSize] = 0x02

		rom := mem.NewROM(bytes)
		ctr, _ := NewROM(rom)

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

		ctr, _ := NewROM(rom)

		got, err := ctr.GetByte(romStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})
}

func Test_romCtr_SetByte(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		bytes[ramSize] = 0x02

		rom := mem.NewROM(bytes)
		ctr, _ := NewROM(rom)

		err := ctr.SetByte(ramEnd, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(ramEnd)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("ROM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)
		ctr, _ := NewROM(rom)

		err := ctr.SetByte(romStart, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(romStart)
		assert.Equal(t, got, byte(0x00))
	})
}

func Test_romCtr_Accepts(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		bytes[ramSize] = 0x02

		rom := mem.NewROM(bytes)
		ctr, _ := NewROM(rom)

		got := ctr.Accepts(ramStart)
		assert.Equal(t, got, true)
	})

	t.Run("ROM address", func(t *testing.T) {
		bytes := make([]byte, romEnd+1)
		rom := mem.NewROM(bytes)
		ctr, _ := NewROM(rom)

		got := ctr.Accepts(0x0000)
		assert.Equal(t, got, true)
	})
}

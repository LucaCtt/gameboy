package cart

import (
	"testing"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func Test_newROMCtr(t *testing.T) {
	t.Run("invalid rom", func(t *testing.T) {
		bytes := make([]byte, 0)
		rom := mem.NewROM(bytes)

		_, err := NewROM(rom)
		assert.Err(t, err, true)
	})

	t.Run("valid rom", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		_, err := NewROM(rom)
		assert.Err(t, err, false)
	})
}

func Test_romCtr_GetByte(t *testing.T) {
	t.Run("RAM addr", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		ctr, _ := NewROM(rom)

		got, err := ctr.GetByte(0xA000)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0xFF))
	})

	t.Run("ROM addr", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		bytes[0x0001] = 0x11
		rom := mem.NewROM(bytes)

		ctr, _ := NewROM(rom)

		got, err := ctr.GetByte(0x0001)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})
}

func Test_romCtr_SetByte(t *testing.T) {
	bytes := make([]byte, 0xFFFF)
	rom := mem.NewROM(bytes)

	ctr, _ := NewROM(rom)

	err := ctr.SetByte(0x0001, 0x11)
	assert.Err(t, err, false)

	got, _ := ctr.GetByte(0x0001)
	assert.Equal(t, got, byte(0x00))
}

func Test_romCtr_Accepts(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		ctr, _ := NewROM(rom)

		got := ctr.Accepts(0xA000)
		assert.Equal(t, got, true)
	})

	t.Run("ROM address", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		rom := mem.NewROM(bytes)

		ctr, _ := NewROM(rom)

		got := ctr.Accepts(0x0000)
		assert.Equal(t, got, true)
	})
}

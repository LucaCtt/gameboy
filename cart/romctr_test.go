package cart

import (
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func Test_newROMCtr(t *testing.T) {
	t.Run("ROM size is too small", func(t *testing.T) {
		_, err := NewROMCtr(make([]byte, 0), make([]byte, 0))
		assert.Err(t, err, true)
	})

	t.Run("ROM size is big enough", func(t *testing.T) {
		_, err := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, 0))
		assert.Err(t, err, false)
	})
}

func Test_ROMCtr_GetByte(t *testing.T) {
	t.Run("RAM addr", func(t *testing.T) {
		ctr, _ := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, ramBankSize))

		// The ram is created by the controller, so it can be accessed
		// only by using SetByte.
		err := ctr.SetByte(romCtrRAMStart, 0x10)
		assert.Err(t, err, false)

		got, err := ctr.GetByte(romCtrRAMStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x10))
	})

	t.Run("ROM addr", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		bytes[romCtrROMEnd] = 0x11

		ctr, _ := NewROMCtr(bytes, make([]byte, 0))

		got, err := ctr.GetByte(romCtrROMEnd)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("Invalid addr", func(t *testing.T) {
		ctr, _ := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, 0))

		_, err := ctr.GetByte(romCtrROMEnd + 1)
		assert.Err(t, err, true)
	})
}

func Test_ROMCtr_SetByte(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		ctr, _ := NewROMCtr(bytes, make([]byte, ramBankSize))

		err := ctr.SetByte(romCtrRAMEnd, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(romCtrRAMEnd)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("ROM address", func(t *testing.T) {
		ctr, _ := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, 0))

		err := ctr.SetByte(romCtrROMEnd, 0x11)
		assert.Err(t, err, false)

		got, _ := ctr.GetByte(romCtrROMEnd)
		assert.Equal(t, got, byte(0x00))
	})

	t.Run("Invalid address", func(t *testing.T) {
		ctr, _ := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, 0))

		err := ctr.SetByte(romCtrROMEnd+1, 0x11)
		assert.Err(t, err, true)
	})
}

func Test_ROMCtr_Accepts(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		ctr, _ := NewROMCtr(bytes, make([]byte, ramBankSize))

		got := ctr.Accepts(romCtrRAMStart)
		assert.Equal(t, got, true)

		got = ctr.Accepts(romCtrRAMEnd)
		assert.Equal(t, got, true)
	})

	t.Run("ROM address", func(t *testing.T) {
		ctr, _ := NewROMCtr(make([]byte, 2*romBankSize), make([]byte, 0))

		got := ctr.Accepts(romCtrROMEnd)
		assert.Equal(t, got, true)

		got = ctr.Accepts(romCtrROMEnd)
		assert.Equal(t, got, true)
	})
}

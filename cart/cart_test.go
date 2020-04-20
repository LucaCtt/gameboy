package cart

import (
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestNewCart(t *testing.T) {
	t.Run("invalid rom", func(t *testing.T) {
		_, err := NewCart(make([]byte, 0))
		assert.Err(t, err, true)
	})

	t.Run("invalid controller", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		bytes[cartTypeFlag] = 0xFF

		_, err := NewCart(bytes)
		assert.Err(t, err, true)
	})

	t.Run("valid controller", func(t *testing.T) {
		_, err := NewCart(make([]byte, romCtrROMEnd+1))
		assert.Err(t, err, false)
	})
}

func TestCart_GetByte(t *testing.T) {
	t.Run("valid addr", func(t *testing.T) {
		bytes := make([]byte, romCtrROMEnd+1)
		bytes[0x0000] = 0x11

		r, _ := NewCart(bytes)

		got, err := r.GetByte(0x0000)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("invalid addr", func(t *testing.T) {
		bytes := make([]byte, romCtrROMEnd+1)

		r, _ := NewCart(bytes)

		_, err := r.GetByte(0xFFFF)
		assert.Err(t, err, true)
	})
}

func TestCart_SetByte(t *testing.T) {
	t.Run("valid addr", func(t *testing.T) {
		bytes := make([]byte, romCtrROMEnd+1)
		bytes[ramSizeFlag] = valueRAMBank1

		r, _ := NewCart(bytes)

		err := r.SetByte(romCtrRAMStart, 0x11)
		assert.Err(t, err, false)

		got, err := r.GetByte(romCtrRAMStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("invalid addr", func(t *testing.T) {
		bytes := make([]byte, romCtrROMEnd+1)

		r, _ := NewCart(bytes)

		err := r.SetByte(0xFFFF, 0x11)
		assert.Err(t, err, true)
	})
}

func TestCart_Accepts(t *testing.T) {
	bytes := make([]byte, romCtrROMEnd+1)
	r, _ := NewCart(bytes)

	got := r.Accepts(romCtrROMStart)
	assert.Equal(t, got, true)
}

func TestCart_Title(t *testing.T) {
	bytes := make([]byte, romCtrROMEnd+1)
	copyAt([]byte{0x54, 0x45, 0x53, 0x54}, bytes, titleStart)

	r, _ := NewCart(bytes)
	assert.Equal(t, r.Title(), "TEST")
}

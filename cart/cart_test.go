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
		bytes[cartType] = 0xFF

		_, err := NewCart(bytes)
		assert.Err(t, err, true)
	})

	t.Run("valid controller", func(t *testing.T) {
		_, err := NewCart(make([]byte, romCtrROMEnd+1))
		assert.Err(t, err, false)
	})
}

func TestReader_Title(t *testing.T) {
	bytes := make([]byte, romCtrROMEnd+1)
	copyAt([]byte{0x54, 0x45, 0x53, 0x54}, bytes, titleStart)

	r, _ := NewCart(bytes)
	assert.Equal(t, r.Title(), "TEST")
}

func copyAt(src, dst []byte, off uint16) {
	for i, b := range src {
		dst[off+uint16(i)] = b
	}
}

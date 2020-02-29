package cart

import (
	"testing"

	"github.com/lucactt/gameboy/gameboy/mem"
	"github.com/lucactt/gameboy/util/assert"
)

func TestNewCart(t *testing.T) {
	t.Run("invalid rom", func(t *testing.T) {
		bytes := make([]byte, 0)
		rom := mem.NewROM(bytes)

		_, err := NewCart(rom)
		assert.Err(t, err, true)
	})

	t.Run("invalid controller", func(t *testing.T) {
		bytes := make([]byte, 0xFFFF)
		bytes[cartType] = 0xFF
		rom := mem.NewROM(bytes)

		_, err := NewCart(rom)
		assert.Err(t, err, true)
	})
}

func TestReader_Title(t *testing.T) {
	bytes := make([]byte, 0xFFFF)
	copyAt([]byte{0x54, 0x45, 0x53, 0x54}, bytes, titleStart)
	rom := mem.NewROM(bytes)

	r, _ := NewCart(rom)
	assert.Equal(t, r.Title(), "TEST")
}

func copyAt(src, dst []byte, off uint16) {
	for i, b := range src {
		dst[off+uint16(i)] = b
	}
}

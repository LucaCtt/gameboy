package cart

import (
	"testing"

	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/assert"
)

func TestNewCart(t *testing.T) {
	bytes := make([]byte, 0)

	rom := memory.NewROM(bytes)

	_, err := NewCart(rom)
	assert.Err(t, err, true)
}

func TestReader_Title(t *testing.T) {
	bytes := make([]byte, headerEnd+1)
	copyAt([]byte{0x54, 0x45, 0x53, 0x54}, bytes, titleStart)

	rom := memory.NewROM(bytes)

	r, _ := NewCart(rom)
	assert.Equal(t, r.Title(), "TEST")
}

func copyAt(src, dst []byte, off uint16) {
	for i, b := range src {
		dst[off+uint16(i)] = b
	}
}
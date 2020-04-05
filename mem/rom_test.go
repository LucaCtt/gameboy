package mem

import (
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestROM_Accepts(t *testing.T) {
	tests := []struct {
		name string
		addr uint16
		want bool
	}{
		{"first byte", 0x0000, true},
		{"last byte", 0x0FFF, true},
		{"upper bound", 0x1000, false},
		{"valid byte", 0x0010, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := make([]byte, 0x1000)
			mem := NewROM(content)

			got := mem.Accepts(tt.addr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestROM_GetByte(t *testing.T) {
	t.Run("inside mem", func(t *testing.T) {
		content := make([]byte, 0x1000)
		content[0x0001] = 0x11
		mem := NewROM(content)

		got, err := mem.GetByte(0x0001)

		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("outside mem", func(t *testing.T) {
		content := make([]byte, 0x1000)
		mem := NewROM(content)

		got, err := mem.GetByte(0x1001)

		assert.Err(t, err, true)
		assert.Equal(t, got, byte(0x00))
	})
}

func TestROM_SetByte(t *testing.T) {
	t.Run("inside mem", func(t *testing.T) {
		content := make([]byte, 0x1000)
		mem := NewROM(content)

		err := mem.SetByte(0x0001, 0x11)
		got, err := mem.GetByte(0x0001)

		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x00))
	})

	t.Run("outside mem", func(t *testing.T) {
		content := make([]byte, 0x1000)
		mem := NewROM(content)

		err := mem.SetByte(0x1001, 0x11)
		got, err := mem.GetByte(0x1001)

		assert.Err(t, err, true)
		assert.Equal(t, got, byte(0x00))
	})
}

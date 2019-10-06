package memory

import (
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestROM_Accepts(t *testing.T) {
	tests := []struct {
		name string
		addr uint16
		want bool
	}{
		{"first byte", 0x0001, true},
		{"last byte", 0x0FFF, true},
		{"zero", 0x0000, false},
		{"upper bound", 0x1000, false},
		{"valid byte", 0x0010, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := make(map[uint16]byte)
			mem := NewROM(0x0001, 0x1000, content)

			got := mem.Accepts(tt.addr)
			util.AssertEqual(t, got, tt.want)
		})
	}
}

func TestROM_GetByte(t *testing.T) {
	t.Run("inside space", func(t *testing.T) {
		content := make(map[uint16]byte)
		content[0x0001] = 0x11
		mem := NewROM(0x0001, 0x1000, content)

		got, err := mem.GetByte(0x0001)

		util.AssertErr(t, err, false)
		util.AssertEqual(t, got, byte(0x11))
	})

	t.Run("outside space", func(t *testing.T) {
		content := make(map[uint16]byte)
		mem := NewROM(0x0001, 0x1000, content)

		got, err := mem.GetByte(0x1001)

		util.AssertErr(t, err, true)
		util.AssertEqual(t, got, byte(0x00))
	})
}

func TestROM_SetByte(t *testing.T) {
	t.Run("inside space", func(t *testing.T) {
		content := make(map[uint16]byte)
		mem := NewROM(0x0001, 0x1000, content)

		err := mem.SetByte(0x0001, 0x11)
		got, err := mem.GetByte(0x0001)

		util.AssertErr(t, err, false)
		util.AssertEqual(t, got, byte(0x00))
	})

	t.Run("outside space", func(t *testing.T) {
		content := make(map[uint16]byte)
		mem := NewROM(0x0001, 0x1000, content)

		err := mem.SetByte(0x1001, 0x11)
		got, err := mem.GetByte(0x1001)

		util.AssertErr(t, err, true)
		util.AssertEqual(t, got, byte(0x00))
	})
}

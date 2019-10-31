package memory

import (
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestNullMem_GetByte(t *testing.T) {
	t.Run("inside space", func(t *testing.T) {
		mem := NewNullMem(0x0001, 0x1000)

		got, err := mem.GetByte(0x0001)

		util.AssertErr(t, err, false)
		util.AssertEqual(t, got, byte(0x00))
	})

	t.Run("outside space", func(t *testing.T) {
		mem := NewNullMem(0x0001, 0x1000)

		got, err := mem.GetByte(0x1001)

		util.AssertErr(t, err, true)
		util.AssertEqual(t, got, byte(0x00))
	})
}

func TestNullMem_SetByte(t *testing.T) {
	t.Run("inside space", func(t *testing.T) {
		mem := NewNullMem(0x0001, 0x1000)

		err := mem.SetByte(0x0001, 0x11)
		got, err := mem.GetByte(0x0001)

		util.AssertErr(t, err, false)
		util.AssertEqual(t, got, byte(0x00))
	})

	t.Run("outside space", func(t *testing.T) {
		mem := NewNullMem(0x0001, 0x1000)

		err := mem.SetByte(0x1001, 0x11)
		got, err := mem.GetByte(0x1001)

		util.AssertErr(t, err, true)
		util.AssertEqual(t, got, byte(0x00))
	})
}

func TestNullMem_Accepts(t *testing.T) {
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
			mem := NewNullMem(0x0001, 0x1000)

			got := mem.Accepts(tt.addr)
			util.AssertEqual(t, got, tt.want)
		})
	}
}

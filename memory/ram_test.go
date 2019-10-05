package memory

import (
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestRAM_Accepts(t *testing.T) {
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
			mem := NewRAM(0x0001, 0x1000)

			got := mem.Accepts(tt.addr)
			util.AssertEqual(t, got, tt.want)
		})
	}
}

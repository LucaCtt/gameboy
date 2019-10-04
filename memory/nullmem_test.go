package memory

import (
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestNullMemory_GetByte(t *testing.T) {
	mem := &NullMemory{}

	got, err := mem.GetByte(0x0001)

	util.AssertErr(t, err, false)
	util.AssertEqual(t, got, byte(0x00))
}

func TestNullMemory_SetByte(t *testing.T) {
	mem := &NullMemory{}

	err := mem.SetByte(0x0001, 0x11)
	got, err := mem.GetByte(0x0001)

	util.AssertErr(t, err, false)
	util.AssertEqual(t, got, byte(0x00))
}

func TestNullMemory_Accepts(t *testing.T) {
	mem := &NullMemory{Start: 0x0000, End: 0x1000}

	got := mem.Accepts(0x0001)

	util.AssertEqual(t, got, true)
}

package cartridge

import (
	"testing"

	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/assert"
)

func TestCartridge_Title(t *testing.T) {
	bytes := make([]byte, hEnd+1)
	bytes[titleStart] = 0x54
	bytes[titleStart+1] = 0x45
	bytes[titleStart+2] = 0x53
	bytes[titleStart+3] = 0x54

	rom := memory.NewROM(bytes)

	r, err := NewReader(rom)
	assert.Err(t, err, false)
	assert.Equal(t, r.Title(), "TEST")
}

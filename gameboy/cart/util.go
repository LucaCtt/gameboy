package cart

import (
	"bytes"
	"fmt"

	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/errors"
)

// getByte returns the byte found at the given address
// in the given rom. It will panic if the address is invalid.
func getByte(rom *memory.ROM, addr uint16) byte {
	b, err := rom.GetByte(addr)

	if err != nil {
		panic(errors.E(fmt.Sprintf("read address %d failed", addr), err, errors.Cartridge))
	}

	return b
}

// getBytes returns a slice of the bytes found at the given address
// range in the given rom. It will panic if the addresses are invalid.
func getBytes(rom *memory.ROM, start, end uint16) []byte {
	result := make([]byte, end-start+1)

	for i := uint16(0); i <= end-start; i++ {
		result[i] = getByte(rom, start+i)
	}

	return result
}

func getString(rom *memory.ROM, start, end uint16) string {
	return string(bytes.Trim(getBytes(rom, start, end), "\x00"))
}

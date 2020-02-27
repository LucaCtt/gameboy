package cart

import (
	"bytes"
	"fmt"

	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/errors"
)

// getByte returns the byte found at the given address
// in the given memory. It will panic if the address is invalid.
func getByte(mem memory.Memory, addr uint16) byte {
	b, err := mem.GetByte(addr)

	if err != nil {
		panic(errors.E(fmt.Sprintf("read address %d failed", addr), err, errors.Cartridge))
	}

	return b
}

// setByte sets the byte at the given address
// in the given memory. It will panic if the address is invalid.
func setByte(mem memory.Memory, addr uint16, value byte) {
	err := mem.SetByte(addr, value)

	if err != nil {
		panic(errors.E(fmt.Sprintf("write to address %d failed", addr), err, errors.Cartridge))
	}
}

// getBytes returns a slice of the bytes found at the given address
// range in the given memory. It will panic if the addresses are invalid.
func getBytes(mem memory.Memory, start, end uint16) []byte {
	result := make([]byte, end-start+1)

	for i := uint16(0); i <= end-start; i++ {
		result[i] = getByte(mem, start+i)
	}

	return result
}

func getString(mem memory.Memory, start, end uint16) string {
	return string(bytes.Trim(getBytes(mem, start, end), "\x00"))
}

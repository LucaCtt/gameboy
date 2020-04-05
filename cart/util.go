package cart

import (
	"bytes"
	"fmt"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// getByte returns the byte found at the given address
// in the given memory. It will panic if the address is invalid.
func getByte(mem mem.Mem, addr uint16) byte {
	b, err := mem.GetByte(addr)

	if err != nil {
		panic(errors.E(fmt.Sprintf("read address %d failed", addr), err, errors.Cartridge))
	}

	return b
}

// setByte sets the byte at the given address
// in the given memory. It will panic if the address is invalid.
func setByte(mem mem.Mem, addr uint16, value byte) {
	err := mem.SetByte(addr, value)

	if err != nil {
		panic(errors.E(fmt.Sprintf("write to address %d failed", addr), err, errors.Cartridge))
	}
}

// getBytes returns a slice of the bytes found at the given address
// range in the given memory. It will panic if the addresses are invalid.
func getBytes(mem mem.Mem, start, end uint16) []byte {
	result := make([]byte, end-start+1)

	for i := uint16(0); i <= end-start; i++ {
		result[i] = getByte(mem, start+i)
	}

	return result
}

func getString(mem mem.Mem, start, end uint16) string {
	return string(bytes.Trim(getBytes(mem, start, end), "\x00"))
}

// romBanks reads the number of ROM banks. Each bank is 32KB.
func romBanks(rom *mem.ROM) int {
	return 2 * (int(getByte(rom, romSize)) ^ 2)
}

func ramBanks(rom *mem.ROM) int {
	switch getByte(rom, ramSize) {
	case 0x02:
		return 1
	case 0x03:
		return 4
	case 0x04:
		return 16
	case 0x05:
		return 8
	default:
		return 0
	}
}

func hasRAM(rom *mem.ROM) bool {
	return ramBanks(rom) != 0
}

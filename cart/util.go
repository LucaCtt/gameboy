package cart

import (
	"bytes"
	"io/ioutil"

	"github.com/lucactt/gameboy/util/errors"
)

// Open reads a file and creates a new cartridge
// with its content.
func Open(p string) (*Cart, error) {
	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.E("read cartridge file failed", err, errors.Cart)
	}

	return NewCart(bytes)
}

// controller wraps a rom with the controller specified by the cart type flag.
func controller(rom []byte, ram []byte) (Controller, error) {
	t := rom[cartType]

	switch {
	case t == 0x00 || t == 0x08 || t == 0x09:
		return NewROMCtr(rom, ram)
	default:
		return nil, errors.E("unsupported cartridge type", errors.Cart)
	}
}

// ramBanks reads the number of RAM banks.
func ramBanks(rom []byte) int {
	switch rom[ramSize] {
	case ramBank1:
		return 1
	case ramBank4:
		return 4
	case ramBank16:
		return 16
	case ramBank8:
		return 8
	default:
		return 0
	}
}

// getString builds a string using the sequence of bytes between two memory addresses,
// trimming any 0x00 byte.
func getString(mem []byte, start, end uint16) string {
	return string(bytes.Trim(mem[start:end], "\x00"))
}

// copyAt copies a source slice into a destination slice, starting at the given address.
func copyAt(src, dst []byte, off uint16) {
	for i, b := range src {
		dst[off+uint16(i)] = b
	}
}

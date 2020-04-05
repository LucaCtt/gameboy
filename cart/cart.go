package cart

import (
	"io/ioutil"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// Cartridge header addresses
const (
	titleStart uint16 = 0x0134
	titleEnd   uint16 = 0x0143
	cartType   uint16 = 0x0147
	romSize    uint16 = 0x0148
	ramSize    uint16 = 0x0149
	headerEnd  uint16 = 0x014F
)

// Controller represents the type of memory bank controller used
// in the cartridge.
type Controller interface {
	mem.Mem
}

// Cart represents a Gameboy cartridge.
type Cart struct {
	title string
	ctr   Controller
}

// NewCart creates a new cartridge from the given ROM.
// It will return an error if the ROM is an invalid cartridge.
func NewCart(rom *mem.ROM) (*Cart, error) {
	if !rom.Accepts(headerEnd) {
		return nil, errors.E("rom size insufficient to contain header", errors.Cartridge)
	}

	title := getString(rom, titleStart, titleEnd)

	ctr, err := controller(rom)
	if err != nil {
		return nil, errors.E("get controller failed", err, errors.Cartridge)
	}

	return &Cart{title, ctr}, nil
}

// Title returns the title of the cartridge.
func (c *Cart) Title() string {
	return c.title
}

// GetByte returns the byte at the given address.
// If the address is not valid, an
// error will be returned.
func (c *Cart) GetByte(addr uint16) (byte, error) {
	b, err := c.ctr.GetByte(addr)
	if err != nil {
		return 0, errors.E("get byte from cartridge failed", err, errors.Cartridge)
	}
	return b, nil
}

// SetByte writes the given value at the given address.
// It will return an error if the address is invalid.
func (c *Cart) SetByte(addr uint16, value byte) error {
	if err := c.ctr.SetByte(addr, value); err != nil {
		return errors.E("get byte from cartridge failed", err, errors.Cartridge)
	}
	return nil
}

// Accepts checks if an address is included in the cartridge.
func (c *Cart) Accepts(addr uint16) bool {
	return c.ctr.Accepts(addr)
}

// Open reads a file and creates a new cartridge
// with its content.
func Open(p string) (*Cart, error) {
	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.E("read cartridge file failed", err, errors.Cartridge)
	}

	return NewCart(mem.NewROM(bytes))
}

func controller(rom *mem.ROM) (Controller, error) {
	t := getByte(rom, cartType)

	switch {
	case t == 0x00 || t == 0x08 || t == 0x09:
		return NewROM(rom)
	default:
		return nil, errors.E("unsupported cartridge type", errors.Cartridge)
	}
}

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

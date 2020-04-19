package cart

import (
	"io/ioutil"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// Cartridge header addresses.
const (
	titleStart uint16 = 0x0134
	titleEnd   uint16 = 0x0143
	cartType   uint16 = 0x0147
	romSize    uint16 = 0x0148
	ramSize    uint16 = 0x0149
	headerEnd  uint16 = 0x014F
)

// Byte values used to identify the number of RAM banks.
// Note that this does not include the values for the ROM banks
// because they can be calculated using the formula:
// rom_bank = 2 * (GetByte(0x0148) ^ 2)
const (
	ramBank1  byte = 0x02
	ramBank4  byte = 0x03
	ramBank16 byte = 0x04
	ramBank8  byte = 0x05

	romBankSize int = 16384
	ramBankSize int = 8192
)

// Controller represents the memory bank controller used
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
// It will return an error if the ROM is an invalid or unsupported cartridge.
func NewCart(rom []byte) (*Cart, error) {
	if len(rom) < int(headerEnd) {
		return nil, errors.E("rom size insufficient to contain header", errors.Cart)
	}

	title := getString(rom, titleStart, titleEnd)

	ctr, err := controller(rom)
	if err != nil {
		return nil, errors.E("create controller failed", err, errors.Cart)
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
		return 0, errors.E("get byte from cartridge failed", err, errors.Cart)
	}
	return b, nil
}

// SetByte writes the given value at the given address.
// It will return an error if the address is invalid.
func (c *Cart) SetByte(addr uint16, value byte) error {
	if err := c.ctr.SetByte(addr, value); err != nil {
		return errors.E("get byte from cartridge failed", err, errors.Cart)
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
		return nil, errors.E("read cartridge file failed", err, errors.Cart)
	}

	return NewCart(bytes)
}

// controller wraps a rom with the controller specified by the cart type flag.
func controller(rom []byte) (Controller, error) {
	t := rom[cartType]

	switch {
	case t == 0x00 || t == 0x08 || t == 0x09:
		return NewROMCtr(rom, 0)
	default:
		return nil, errors.E("unsupported cartridge type", errors.Cart)
	}
}

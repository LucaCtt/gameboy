package cart

import (
	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// Addresses of the info contained in the header.
const (
	titleStart   uint16 = 0x0134
	titleEnd     uint16 = 0x0143
	cartTypeFlag uint16 = 0x0147
	romSizeFlag  uint16 = 0x0148
	ramSizeFlag  uint16 = 0x0149
	headerEnd    uint16 = 0x014F
)

// Byte values used to identify the number of RAM banks.
const (
	valueRAMBank1  byte = 0x02
	valueRAMBank4  byte = 0x03
	valueRAMBank16 byte = 0x04
	valueRAMBank8  byte = 0x05
)

// Size of the ROM and RAM banks.
const (
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

	ram := make([]byte, ramBanks(rom))
	ctr, err := controller(rom, ram)
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

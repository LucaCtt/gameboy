package cart

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// Memory address ranges.
const (
	romCtrROMStart uint16 = 0x0000
	romCtrROMEnd   uint16 = 0x7FFF
	romCtrRAMStart uint16 = 0xA000
	romCtrRAMEnd   uint16 = 0xBFFF
)

// ROMCtr represents a ROM controller.
//
// It maps a ROM to 0x0000 to 0x7FFF, and can optionally
// map a RAM to 0xA000-0xBFFF.
type ROMCtr struct {
	rom []byte
	ram []byte
}

// NewROMCtr creates a new ROM controller from the given ROM and RAM banks number.
//
// The ROM must be large enough to contain at least two banks.
// The number of RAM banks can be 0.
func NewROMCtr(rom []byte, ramBanks int) (*ROMCtr, error) {
	if len(rom) < 2*romBankSize {
		return nil, errors.E("rom size insufficient: must contain at least two banks", errors.Cart)
	}

	ram := make([]byte, int(ramBanks)*ramBankSize)
	return &ROMCtr{rom: rom, ram: ram}, nil
}

// GetByte returns the byte at the given address, which
// can be read from the ROM or from the RAM, if it exists.
func (ctr *ROMCtr) GetByte(addr uint16) (byte, error) {
	if addr <= romCtrROMEnd {
		return ctr.rom[addr], nil
	} else if addr >= romCtrRAMStart && addr <= romCtrRAMEnd {
		return ctr.ram[addr-romCtrRAMStart], nil
	}

	return 0, errors.E(fmt.Sprintf("ROM controller doesn't accept addr %d", addr),
		errors.CodeOutOfRange,
		errors.Cart)
}

// SetByte does nothing if the addr points to the ROM,
// or sets the byte to the given value if it points to RAM.
func (ctr *ROMCtr) SetByte(addr uint16, value byte) error {
	if addr <= romCtrROMEnd {
		return nil
	} else if addr >= romCtrRAMStart && addr <= romCtrRAMEnd {
		ctr.ram[addr-romCtrRAMStart] = value
		return nil
	} else {
		return errors.E(fmt.Sprintf("ROM controller doesn't accept addr %d", addr),
			errors.CodeOutOfRange,
			errors.Cart)
	}
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *ROMCtr) Accepts(addr uint16) bool {
	return (addr >= romCtrROMStart && addr <= romCtrROMEnd) ||
		(len(ctr.ram) > 0 && addr >= romCtrRAMStart && addr <= romCtrRAMEnd)
}

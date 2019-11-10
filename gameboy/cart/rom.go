package cart

import (
	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/errors"
)

// ROM represents a ROM controller.
// It maps a rom to 0x0000 to 0x7FFF directly.
//
// At the moment it does not support a RAM,
// so any reads from 0xA000-0xBFFF will return 0xFF
// and writes will have no effect.
type ROM struct {
	rom *memory.ROM
}

// NewROM creates a new ROM controller with the given rom.
// The rom must at least include addresses until 0x7FFF, otherwise
// an error will be returned.
func NewROM(rom *memory.ROM) (*ROM, error) {
	if !rom.Accepts(0x7FFF) {
		return nil, errors.E("rom size insufficient", errors.Cartridge)
	}

	return &ROM{rom}, nil
}

// GetByte returns the byte at the given address, or
// 0xFF if the address points to the RAM.
func (ctr *ROM) GetByte(addr uint16) (byte, error) {
	if isRAM(addr) {
		return 0xFF, nil
	}

	return getByte(ctr.rom, addr), nil
}

// SetByte does nothing.
func (ctr *ROM) SetByte(addr uint16, value byte) error {
	return nil
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *ROM) Accepts(addr uint16) bool {
	return (addr >= 0x0000 && addr < 0x8000) || isRAM(addr)
}

func isRAM(addr uint16) bool {
	return addr >= 0xA000 && addr < 0xC000
}

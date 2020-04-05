package cart

import (
	"fmt"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

const (
	romStart uint16 = 0x0000
	romEnd   uint16 = 0x7FFF
	ramStart uint16 = 0xA000
	ramEnd   uint16 = 0xBFFF
)

// ROM represents a ROM controller.
// It maps a rom to 0x0000 to 0x7FFF directly, and can optionally
// map a ram to 0xA000-0xBFFF.
type ROM struct {
	mem    mem.Mem
	hasRAM bool
}

// NewROM creates a new ROM controller with the given rom.
// The rom must at least include addresses until 0x7FFF, otherwise
// an error will be returned.
//
// If the ram banks flag (0x0149) in the passed rom is not zero,
// an 8MB ram will also be created, and it will accept addresses between 0xA000-0xBFFF.
func NewROM(rom *mem.ROM) (*ROM, error) {
	if !rom.Accepts(romEnd) {
		return nil, errors.E("rom size insufficient", errors.Cartridge)
	}

	if hasRAM(rom) {
		mmu := &mem.MMU{}
		mmu.AddMem(romStart, rom)
		mmu.AddMem(ramStart, mem.NewRAM(ramEnd+1-ramStart))

		return &ROM{mmu, true}, nil
	}

	return &ROM{rom, false}, nil
}

// GetByte returns the byte at the given address.
func (ctr *ROM) GetByte(addr uint16) (byte, error) {
	if !ctr.Accepts(addr) {
		return 0, errors.E(fmt.Sprintf("ROM cart can't get addr %d", addr),
			errors.CodeOutOfRange,
			errors.Cartridge)
	}

	return getByte(ctr.mem, addr), nil
}

// SetByte does nothing if the addr points to the rom,
// or sets the byte to the value if it is in the ram.
func (ctr *ROM) SetByte(addr uint16, value byte) error {
	if !ctr.Accepts(addr) {
		return errors.E(fmt.Sprintf("ROM cart can't set addr %d", addr),
			errors.CodeOutOfRange,
			errors.Cartridge)
	}

	return ctr.mem.SetByte(addr, value)
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *ROM) Accepts(addr uint16) bool {
	a := (addr >= 0x0000 && addr < 0x8000)
	b := (ctr.hasRAM && addr >= 0xA000 && addr < 0xC000)
	return a || b
}

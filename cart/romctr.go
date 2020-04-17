package cart

import (
	"fmt"

	"github.com/lucactt/gameboy/mem"
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
	mem    mem.Mem
	hasRAM bool
}

// NewROMCtr creates a new ROM controller with the given ROM.
// The RAM must at least include addresses until 0x7FFF, otherwise
// an error will be returned.
//
// If the RAM banks flag (0x0149) in the given rom is not zero,
// an 8MB RAM will also be created, which will accept addresses between 0xA000-0xBFFF.
//
// Note that
func NewROMCtr(rom *mem.ROM) (*ROMCtr, error) {
	if !rom.Accepts(romCtrROMEnd) {
		return nil, errors.E("ROM size insufficient", errors.Cartridge)
	}

	if rom.Accepts(romCtrRAMStart) {
		return nil, errors.E("ROM size is too big: it covers the RAM addresses", errors.Cartridge)
	}

	if ramBanks(rom) != 0 {
		mmu := &mem.MMU{}
		mmu.AddMem(romCtrROMStart, rom)
		mmu.AddMem(romCtrRAMStart, mem.NewRAM(romCtrRAMEnd-romCtrRAMStart+1))

		return &ROMCtr{mmu, true}, nil
	}

	return &ROMCtr{rom, false}, nil
}

// GetByte returns the byte at the given address, which
// can be read from the ROM or from the RAM, if it exists.
func (ctr *ROMCtr) GetByte(addr uint16) (byte, error) {
	if !ctr.Accepts(addr) {
		return 0, errors.E(fmt.Sprintf("ROM controller doesn't accept addr %d", addr),
			errors.CodeOutOfRange,
			errors.Cartridge)
	}

	return getByte(ctr.mem, addr), nil
}

// SetByte does nothing if the addr points to the ROM,
// or sets the byte to the given value if it points to RAM.
func (ctr *ROMCtr) SetByte(addr uint16, value byte) error {
	if !ctr.Accepts(addr) {
		return errors.E(fmt.Sprintf("ROM controller doesn't accept addr %d", addr),
			errors.CodeOutOfRange,
			errors.Cartridge)
	}

	return ctr.mem.SetByte(addr, value)
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *ROMCtr) Accepts(addr uint16) bool {
	return (addr >= romCtrROMStart && addr <= romCtrROMEnd) ||
		(ctr.hasRAM && addr >= romCtrRAMStart && addr <= romCtrRAMEnd)
}

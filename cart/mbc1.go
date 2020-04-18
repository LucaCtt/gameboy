package cart

import (
	"fmt"

	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// Memory address ranges
const (
	mbc1ROMBank0Start  uint16 = 0x0000
	mbc1ROMBank0End    uint16 = 0x3FFF
	mbc1SwitchROMStart uint16 = 0x4000
	mbc1SwitchROMEnd   uint16 = 0x7FFF
	mbc1SwitchRAMStart uint16 = 0xA000
	mbc1SwitchRAMEnd   uint16 = 0xBFFF
	mbc1RAMEnableStart uint16 = 0x0000
	mbc1RAMEnableEnd   uint16 = 0x1FFF
	mbc1ROMBankStart   uint16 = 0x2000
	mbc1ROMBankEnd     uint16 = 0x3FFF
	mbc1RAMBankStart   uint16 = 0x4000
	mbc1RAMBankEnd     uint16 = 0x5FFF
	mbc1ModeStart      uint16 = 0x6000
	mbc1ModeEnd        uint16 = 0x7FFF

	mbc1EnableRAMValue byte   = 0x0A
	mbc1ROMBankSize    uint16 = mbc1SwitchROMEnd + 1 - mbc1SwitchROMStart
	mbc1RAMBankSize    uint16 = mbc1SwitchRAMEnd + 1 - mbc1SwitchRAMStart
)

// MBC1 implements an MBC1 cartridge controller.
type MBC1 struct {
	mem          mem.Mem
	romBank      byte
	ramBank      byte
	isROMBanking bool
	isRAMEnabled bool
}

// NewMBC1 creates a new MBC1 controller from the given ROM and RAM.
//
// The ROM must at least be big enough to contain the addresses until 0x7FFF,
// which is the same as saying that it must have at least two banks.
//
// If the RAM is not nil, it must at least be big enough to accept the addresses
// between 0xA000 and 0xBFFF.
func NewMBC1(rom *mem.ROM, ram *mem.RAM) (*MBC1, error) {
	if !rom.Accepts(mbc1SwitchROMEnd) {
		return nil, errors.E("rom size insufficient", errors.Cart)
	}

	if ram == nil {
		return &MBC1{mem: rom}, nil
	}

	if !ram.Accepts(mbc1SwitchRAMEnd - mbc1SwitchRAMStart) {
		return nil, errors.E("ram size insufficient", errors.Cart)
	}

	mmu := &mem.MMU{}
	mmu.AddMem(mbc1ROMBank0Start, rom)
	mmu.AddMem(mbc1SwitchRAMStart, ram)

	return &MBC1{mem: mmu}, nil
}

// GetByte returns the byte at the given address, which
// can be read from the ROM or from the RAM, if it exists.
func (ctr *MBC1) GetByte(addr uint16) (byte, error) {
	if !ctr.Accepts(addr) {
		return 0, errors.E(fmt.Sprintf("mbc1 controller does not accept addr %d", addr))
	}

	// This switch only handles cases where the address is in the switch rom or switch ram.
	// Every other case is handled by the default.
	switch {
	case addr >= mbc1SwitchROMStart && addr <= mbc1SwitchROMEnd:
		// Sum the address corresponding to the bank number to the address
		// given as parameter.
		relAddr := uint16(ctr.romBank)*mbc1ROMBankSize + addr

		res, err := ctr.mem.GetByte(relAddr)
		if err != nil {
			panic(errors.E(
				fmt.Sprintf("get rom addr %d in bank %d of mbc1 failed", relAddr, ctr.romBank),
				err))
		}
		return res, nil

	case addr >= mbc1SwitchRAMStart && addr <= mbc1SwitchRAMEnd:
		if !ctr.isRAMEnabled {
			return 0xFF, nil
		}

		relAddr := uint16(ctr.ramBank)*mbc1RAMBankSize + addr

		res, err := ctr.mem.GetByte(relAddr)
		if err != nil {
			panic(errors.E(
				fmt.Sprintf("get ram addr %d in bank %d of mbc1 failed", relAddr, ctr.ramBank),
				err))
		}
		return res, nil

	default:
		res, err := ctr.mem.GetByte(addr)
		if err != nil {
			panic(errors.E(fmt.Sprintf("mbc1 mem accepts addr %d, but GetByte returns err", addr), err))
		}
		return res, nil
	}
}

// SetByte does nothing if the addr points to the ROM,
// or sets the byte to the given value if it points to RAM.
func (ctr *MBC1) SetByte(addr uint16, value byte) error {
	if !ctr.Accepts(addr) {
		return errors.E(fmt.Sprintf("mbc1 controller does not accept addr %d", addr))
	}

	switch {
	case addr >= mbc1RAMEnableStart && addr <= mbc1RAMEnableEnd:
		ctr.isRAMEnabled = (value == mbc1EnableRAMValue)

	case addr >= mbc1ROMBankStart && addr <= mbc1ROMBankEnd:
		bank := value

		if bank == 0x00 || bank == 0x20 || bank == 0x40 || bank == 0x60 {
			bank += 0x01
		}

		ctr.romBank = bank & 0x1F
	default:
		if err := ctr.mem.SetByte(addr, value); err != nil {
			panic(errors.E(fmt.Sprintf("mbc1 mem accepts addr %d, but SetByte returns err", addr), err))
		}
	}

	return nil
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *MBC1) Accepts(addr uint16) bool {
	return ctr.mem.Accepts(addr)
}

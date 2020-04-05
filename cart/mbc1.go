package cart

import "github.com/lucactt/gameboy/mem"

type MBC1 struct {
	rom                *mem.ROM
	ram                *mem.RAM
	romBanks, ramBanks int
	ramEnabled         bool
	isROMBanking       bool
	romBank, ramBank   int
}

func NewMBC1(rom *mem.ROM, romBanks, ramBanks int) *MBC1 {
	ram := mem.NewRAM(uint16(0x2000 * ramBanks))
	return &MBC1{rom, ram, romBanks, ramBanks, false, true, 1, 1}
}

func (ctr *MBC1) GetByte(addr uint16) (byte, error) {
	if isRAM(addr) {
		return 0xFF, nil
	}

	return getByte(ctr.rom, addr), nil
}

func (ctr *MBC1) SetByte(addr uint16, value byte) error {
	switch {
	case addr >= 0x000 && addr <= 0x1FFF:
		ctr.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr <= 0x3FFF:
		bank := (ctr.romBank & 0x60) | (int(value) & 0x1F)
		ctr.romBank = bank

	case addr >= 0x4000 && addr <= 0x5FFF:
		if ctr.isROMBanking {
			bank := ctr.romBank&0x1F | ((int(value) & 0x03) << 5)
			ctr.romBank = bank
		} else {
			bank := int(value) & 0x03
			ctr.ramBank = bank
		}

	case addr >= 0x6000 && addr <= 0x7FFF:
		ctr.isROMBanking = (value == 0x00)

	case addr >= 0xA000 && addr < 0xC000 && ctr.ramEnabled:
		ramAddr := ctr.ramAddr(addr)
		setByte(ctr.ram, ramAddr, value)
	}
	return nil
}

func (ctr *MBC1) Accepts(addr uint16) bool {
	return (addr >= 0x0000 && addr < 0x8000) ||
		(addr >= 0xA000 && addr < 0xC000)
}

func (ctr *MBC1) ramAddr(addr uint16) uint16 {
	if ctr.isROMBanking {
		return addr - 0xA000
	}

	return uint16(ctr.ramBank%ctr.ramBanks)*0x2000 + (addr - 0xa000)
}

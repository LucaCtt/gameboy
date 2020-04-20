package cart

import (
	"fmt"

	"github.com/lucactt/gameboy/util/errors"
)

// Memory addresses
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

	mbc1EnableRAMValue byte = 0x0A
)

// MBC1 implements an MBC1 cartridge controller.
type MBC1 struct {
	rom          []byte
	ram          []byte
	romBank      byte
	ramBank      byte
	isRAMBanking bool
	isRAMEnabled bool
}

// NewMBC1 creates a new MBC1 controller from the given ROM and RAM banks number.
//
// The ROM must be large enough to contain at least two banks.
// The RAM can have length == 0, but cannot be nil.
func NewMBC1(rom []byte, ram []byte) (*MBC1, error) {
	if rom == nil || ram == nil {
		panic(fmt.Errorf("the rom or ram are nil"))
	}

	if len(rom) < 2*romBankSize {
		return nil, errors.E("rom size insufficient: must contain at least two banks", errors.Cart)
	}

	// The ROM bank is initialized to 0x01 to avoid access to ROM banks 0x00, 0x20, 0x40 and 0x60
	// from the switchable ROM addresses on startup.
	// The SetByte method verifies that the lower two bits of the bank are also != 00 to impose this
	// after startup.
	return &MBC1{rom: rom, ram: ram, romBank: 0x01}, nil
}

// GetByte returns the byte at the given address, which
// can be read from the ROM or from the RAM, if it exists.
func (ctr *MBC1) GetByte(addr uint16) (byte, error) {
	// This also verifies that the RAM cannot be accessed if the controllor has no RAM.
	if !ctr.Accepts(addr) {
		return 0, errors.E(fmt.Sprintf("mbc1 controller does not accept addr %d", addr))
	}

	switch {
	case addr <= mbc1ROMBank0End:
		return ctr.rom[addr], nil

	case addr <= mbc1SwitchROMEnd:
		// The address to read is equal to the selected ROM bank multiplied by the ROM bank size,
		// to which the param address is summed. However, the ROM stored in the controller
		// starts addressing from 0, so the start address of the ROM must be subtracted from the param address.
		relAddr := int(ctr.romBank)*romBankSize + int(addr-mbc1SwitchROMStart)
		return ctr.rom[relAddr], nil

	case addr >= mbc1SwitchRAMStart && addr <= mbc1SwitchRAMEnd:
		if !ctr.isRAMEnabled {
			return 0xFF, nil
		}

		relAddr := int(ctr.ramBank)*ramBankSize + int(addr-mbc1SwitchRAMStart)
		return ctr.ram[relAddr], nil

	default:
		panic(fmt.Errorf("unhandled address %d in mbc1 controller", addr))
	}
}

// SetByte does nothing if the addr points to the ROM,
// or sets the byte to the given value if it points to RAM.
func (ctr *MBC1) SetByte(addr uint16, value byte) error {
	if !ctr.Accepts(addr) {
		return errors.E(fmt.Sprintf("mbc1 controller does not accept addr %d", addr))
	}

	switch {
	case addr <= mbc1RAMEnableEnd:
		ctr.isRAMEnabled = (value == mbc1EnableRAMValue)

	case addr <= mbc1ROMBankEnd:
		bank := value
		if bank == 0x00 {
			bank = 0x01
		}

		// Select only the 5th and 6th bits of the current bank
		bankUpper := (ctr.romBank >> 5) & 0x03
		// Select only the last five bits
		bankLower := bank & 0x1F
		// Set the 5th and 6th bit of the bank, leave the 7th untouched
		ctr.romBank = bankLower + (bankUpper << 5)

	case addr <= mbc1RAMBankEnd:
		if ctr.isRAMBanking {
			ctr.ramBank = value & 0x03
		} else {
			bankUpper := value & 0x03
			bankLower := ctr.romBank & 0x1F
			ctr.romBank = bankLower + (bankUpper << 5)
		}

	case addr <= mbc1ModeEnd:
		ctr.isRAMBanking = (value != 0x00)

		if ctr.isRAMBanking {
			// Clear the upper bits of ROM bank
			ctr.romBank = ctr.romBank & 0x1F
		} else {
			ctr.ramBank = 0
		}

	case addr >= mbc1SwitchRAMStart && addr <= mbc1SwitchRAMEnd:
		relAddr := int(ctr.ramBank)*ramBankSize + int(addr-mbc1SwitchRAMStart)
		ctr.ram[relAddr] = value

	default:
		panic(fmt.Errorf("unhandled address %d in mbc1 controller", addr))
	}

	return nil
}

// Accepts returns true if the address is included in the ROM
// or in the RAM, false otherwise.
func (ctr *MBC1) Accepts(addr uint16) bool {
	if addr >= mbc1SwitchRAMStart && addr <= mbc1SwitchRAMEnd && len(ctr.ram) == 0 {
		return false
	}

	return (addr <= mbc1SwitchROMEnd) || (addr >= mbc1SwitchRAMStart && addr <= mbc1SwitchRAMEnd)
}

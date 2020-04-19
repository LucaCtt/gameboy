package cart

import (
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestNewMBC1(t *testing.T) {
	t.Run("ROM is too small", func(t *testing.T) {
		_, err := NewMBC1(make([]byte, 0), 0)
		assert.Err(t, err, true)
	})

	t.Run("ROM is big enough", func(t *testing.T) {
		_, err := NewMBC1(make([]byte, 2*romBankSize), 0)
		assert.Err(t, err, false)
	})
}

func TestMBC1_GetByte(t *testing.T) {
	t.Run("ROM bank", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		bytes[mbc1ROMBank0End] = 0x11

		ctr, _ := NewMBC1(bytes, 0)

		got, err := ctr.GetByte(mbc1ROMBank0End)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("RAM bank, RAM enabled", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 1)

		// Make sure RAM is enabled
		ctr.SetByte(mbc1RAMEnableStart, mbc1EnableRAMValue)
		ctr.SetByte(mbc1SwitchRAMEnd, 0x11)

		got, err := ctr.GetByte(mbc1SwitchRAMEnd)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("RAM bank, RAM disabled", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 1)

		// Make sure RAM is disabled
		ctr.SetByte(mbc1RAMEnableStart, 0x00)

		got, err := ctr.GetByte(mbc1SwitchRAMEnd)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0xFF))
	})

	t.Run("RAM bank, but no RAM", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 0)

		ctr.SetByte(mbc1SwitchRAMEnd, 0x11)

		_, err := ctr.GetByte(mbc1SwitchRAMEnd)
		assert.Err(t, err, true)
	})
}

func TestMBC1_SetByte(t *testing.T) {
	t.Run("Disable RAM", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 1)

		ctr.SetByte(mbc1RAMEnableStart, 0x00)

		got, err := ctr.GetByte(mbc1SwitchRAMStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0xFF))
	})

	t.Run("Enable RAM", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 1)

		ctr.SetByte(mbc1RAMEnableStart, mbc1EnableRAMValue)
		ctr.SetByte(mbc1SwitchRAMStart, 0x11)

		got, err := ctr.GetByte(mbc1SwitchRAMStart)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("Enable ROM banking", func(t *testing.T) {
		/*
			bytes := make([]byte, 32*romBankSize)
			bytes[32*(romBankSize)-1] = 0x11

			ctr, _ := NewMBC1(bytes, 2)
			ctr.SetByte(mbc1ModeStart, 0x00)
			ctr.SetByte(mbc1RAMBankStart, 0x01)

			got, err := ctr.GetByte(mbc1SwitchROMStart)
			assert.Err(t, err, false)
			assert.Equal(t, got, byte(0x11))*/
	})

	t.Run("ROM banking by default", func(t *testing.T) {
	})

	t.Run("Enable RAM banking", func(t *testing.T) {
	})

	t.Run("ROM banking allows access only to RAM bank 0 ", func(t *testing.T) {
	})

	t.Run("RAM banking allows access only to ROM banks 0x01-0x1F", func(t *testing.T) {
	})

	t.Run("Switch lower 5 bits of ROM bank", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		bytes[2*(romBankSize)-1] = 0x11

		ctr, _ := NewMBC1(bytes, 0)
		ctr.SetByte(mbc1ROMBankStart, 0x01)

		got, err := ctr.GetByte(mbc1SwitchROMEnd)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})

	t.Run("Lower bit of ROM bank converted to 1 if 0", func(t *testing.T) {
		bytes := make([]byte, 2*romBankSize)
		bytes[2*(romBankSize)-1] = 0x11

		ctr, _ := NewMBC1(bytes, 0)
		ctr.SetByte(mbc1ROMBankStart, 0x00)

		got, err := ctr.GetByte(mbc1SwitchROMEnd)
		assert.Err(t, err, false)
		assert.Equal(t, got, byte(0x11))
	})
}

func TestMBC1_Accepts(t *testing.T) {
	t.Run("RAM address", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 1)

		got := ctr.Accepts(mbc1SwitchRAMStart)
		assert.Equal(t, got, true)

		got = ctr.Accepts(mbc1SwitchRAMEnd)
		assert.Equal(t, got, true)
	})

	t.Run("RAM address, but no RAM", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 0)

		got := ctr.Accepts(mbc1SwitchRAMStart)
		assert.Equal(t, got, false)

		got = ctr.Accepts(mbc1SwitchRAMEnd)
		assert.Equal(t, got, false)
	})

	t.Run("ROM address", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 0)

		got := ctr.Accepts(mbc1ROMBank0Start)
		assert.Equal(t, got, true)

		got = ctr.Accepts(mbc1ROMBank0End)
		assert.Equal(t, got, true)
	})

	t.Run("Outside mem", func(t *testing.T) {
		ctr, _ := NewMBC1(make([]byte, 2*romBankSize), 0)

		got := ctr.Accepts(0xFFFF)
		assert.Equal(t, got, false)

		got = ctr.Accepts(mbc1SwitchRAMEnd + 1)
		assert.Equal(t, got, false)
	})
}

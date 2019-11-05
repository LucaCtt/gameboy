package cartridge

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/lucactt/gameboy/gameboy/memory"
	"github.com/lucactt/gameboy/util/errors"
)

// nintendoLogo is the correct Nintendo logo, used for
// comparison with the one in the cartridge.
var nintendoLogo = [48]byte{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

const (
	logoStart   uint16 = 0x0104
	logoEnd     uint16 = 0x0133
	titleStart  uint16 = 0x0134
	titleEnd    uint16 = 0x0143
	gbcFlag     uint16 = 0x0143
	cartType    uint16 = 0x0147
	romSize     uint16 = 0x0148
	ramSize     uint16 = 0x0149
	hCheck      uint16 = 0x014D
	hCheckStart uint16 = 0x0134
	hCheckEnd   uint16 = 0x014C
	hEnd        uint16 = 0x014F
)

// Type represents the type of controller
// in the cartridge.
type Type string

// Controller types.
const (
	MBC1 Type = "MBC1"
	MBC2 Type = "MBC2"
	MBC3 Type = "MBC3"
	MBC5 Type = "MBC5"
	MBC6 Type = "MBC6"
	MBC7 Type = "MBC7"
	ROM  Type = "ROM"
)

// Reader represents a GameBoy cartridge data reader.
//
// If the MBC flag byte is invalid or unsupported,
// the cartridge will be treated as ROM only.
// Currently, only ROM and MBC cartridges are supported.
type Reader struct {
	rom *memory.ROM
}

// NewReader creates a new cartridge reader with the given ROM
//
// The ROM must at least be large enough
// to contain the cartridge header.
func NewReader(rom *memory.ROM) (*Reader, error) {
	if !rom.Accepts(hEnd) {
		return nil, errors.E("rom size insufficient to contain header", errors.Cartridge)
	}

	return &Reader{rom}, nil
}

// ROM returns the ROM read by the Reader.
func (r *Reader) ROM() *memory.ROM {
	return r.rom
}

// Logo returns the value of the Nintendo logo
// contained in the cartridge.
func (r *Reader) Logo() [48]byte {
	logo := [48]byte{}
	copy(logo[:], r.getBytes(logoStart, logoEnd))

	return logo
}

// Title returns the title of the game.
func (r *Reader) Title() string {
	temp := r.getBytes(titleStart, titleEnd)
	return string(bytes.Trim(temp, "\x00"))
}

// IsGBCOnly returns true if the cartridge can only run on the GameBoy Color.
func (r *Reader) IsGBCOnly() bool {
	flag := r.getByte(gbcFlag)
	return flag == 0xC0
}

// Type returns true if the cartridge contains a MBC1.
func (r *Reader) Type() Type {
	t := r.getByte(cartType)
	switch {
	case t >= 0x01 && t <= 0x03:
		return MBC1
	case t == 0x05 || t == 0x06:
		return MBC2
	case t >= 0x0F && t <= 0x13:
		return MBC3
	case t >= 0x19 && t <= 0x1E:
		return MBC5
	case t == 0x20:
		return MBC6
	case t == 0x22:
		return MBC7
	default:
		return ROM
	}
}

// HasBattery returns true if the cartridge containts a battery.
func (r *Reader) HasBattery() bool {
	t := r.getByte(cartType)
	hasBattery := []byte{0x03, 0x06, 0x09, 0x0D, 0x0F, 0x10, 0x13, 0x1B, 0x1E, 0x20, 0x22}

	for _, b := range hasBattery {
		if t == b {
			return true
		}
	}

	return false
}

// HeaderChecksum returns the header checksum contained in the cartridge.
func (r *Reader) HeaderChecksum() byte {
	return r.getByte(hCheck)
}

// IsValid checks the cartridge validity by
// verifying the Nintendo logo, the
// header checksum, and that the cartridge is not GBC only.
//
// This does not verify the global checksum
// in order to emulate the GameBoy behavior accurately.
func (r *Reader) IsValid() bool {
	return !r.IsGBCOnly() &&
		r.Logo() == nintendoLogo &&
		r.computeChecksum() == r.HeaderChecksum()
}

func (r *Reader) getByte(addr uint16) byte {
	b, err := r.rom.GetByte(addr)

	if err != nil {
		panic(errors.E(fmt.Sprintf("read address %d failed", addr), err, errors.Cartridge))
	}

	return b
}

func (r *Reader) getBytes(start, end uint16) []byte {
	result := make([]byte, end-start+1)

	for i := uint16(0); i <= end-start; i++ {
		result[i] = r.getByte(start + i)
	}

	return result
}

func (r *Reader) computeChecksum() byte {
	var sum byte
	for _, b := range r.getBytes(hCheckStart, hCheckEnd) {
		sum -= b - 1
	}
	return sum
}

func readROM(p string) ([]byte, error) {
	if path.Ext(p) != ".gb" {
		return nil, errors.E("invalid cartridge file format", errors.Cartridge)
	}

	return ioutil.ReadFile(p)
}

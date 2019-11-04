package cartridge

import (
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
	manufStart  uint16 = 0x013F
	manufEnd    uint16 = 0x0142
	gbcFlag     uint16 = 0x0143
	newLicStart uint16 = 0x0144
	newLicEnd   uint16 = 0x0145
	sgbFlag     uint16 = 0x0146
	cartType    uint16 = 0x0147
	romSize     uint16 = 0x0148
	ramSize     uint16 = 0x0149
	destCode    uint16 = 0x014A
	oldLic      uint16 = 0x014B
	romVers     uint16 = 0x014C
	hCheck      uint16 = 0x014D
	gCheckStart uint16 = 0x014E
	gCheckEnd   uint16 = 0x014F
)

// Reader represents a GameBoy cartridge data reader.
//
// If the MBC flag byte is invalid or unsupported,
// the cartridge will be treated as ROM only.
// Currently, only ROM and MBC cartridges are supported.
type Reader struct {
	rom *memory.ROM
}

// New creates a new cartridge reader with the ROM
// at the given path.
//
// The ROM file must have a ".gb" file extension and
// its total size must at least be enough
// to contain the cartridge header.
func New(path string) (*Reader, error) {
	data, err := readROM(path)

	if err != nil {
		return nil, errors.E("cannot read ROM file", err, errors.Cartridge)
	}

	rom := memory.NewROM(data)
	if !rom.Accepts(gCheckEnd) {
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

// Title returns the title of the game in uppercase.
func (r *Reader) Title() string {
	return string(r.getBytes(titleStart, titleEnd))
}

// ManufacturerCode returns the manufacturer code in uppercase.
func (r *Reader) ManufacturerCode() string {
	return string(r.getBytes(manufStart, manufEnd))
}

// IsGBCOnly returns true if the cartridge can only run on the GameBoy Color.
func (r *Reader) IsGBCOnly() bool {
	flag := r.getByte(gbcFlag)
	return flag == 0xC0
}

// LicenseeCode returns the licensee code in uppercase.
func (r *Reader) LicenseeCode() string {
	if r.IsSGB() {
		return string(r.getBytes(newLicStart, newLicEnd))
	}
	return string(r.getByte(oldLic))
}

// IsSGB returns true if the cartridge supports SGB functions,
// false otherwise.
func (r *Reader) IsSGB() bool {
	return r.getByte(sgbFlag) == 0x03
}

// IsMBC1 returns true if the cartridge contains a MBC1.
func (r *Reader) IsMBC1() bool {
	t := r.getByte(cartType)
	return t >= 0x01 && t <= 0x03
}

// IsMBC2 returns true if the cartridge contains a MBC2.
func (r *Reader) IsMBC2() bool {
	t := r.getByte(cartType)
	return t == 0x05 || t == 0x06
}

// IsMBC3 returns true if the cartridge contains a MBC3.
func (r *Reader) IsMBC3() bool {
	t := r.getByte(cartType)
	return t >= 0x0F && t <= 0x13
}

// IsMBC5 returns true if the cartridge contains a MBC5.
func (r *Reader) IsMBC5() bool {
	t := r.getByte(cartType)
	return t >= 0x19 && t <= 0x1E
}

// IsMBC6 returns true if the cartridge contains a MBC6.
func (r *Reader) IsMBC6() bool {
	return r.getByte(cartType) == 0x20
}

// IsMBC7 returns true if the cartridge contains a MBC7.
func (r *Reader) IsMBC7() bool {
	return r.getByte(cartType) == 0x22
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

// IsJapanOnly returns true if the cartridge is supposed to be sold
// in Japan, false otherwise.
func (r *Reader) IsJapanOnly() bool {
	return r.getByte(destCode) == 0x00
}

// ROMVersion returns a byte that indicates the version of the ROM.
func (r *Reader) ROMVersion() byte {
	return r.getByte(romVers)
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
		panic(err)
	}

	return b
}

func (r *Reader) getBytes(start, end uint16) []byte {
	result := make([]byte, end-start)

	for i := start; i <= end; i++ {
		result[i] = r.getByte(i)
	}

	return result
}

func (r *Reader) computeChecksum() byte {
	var sum byte
	for _, b := range r.getBytes(titleStart, romVers) {
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

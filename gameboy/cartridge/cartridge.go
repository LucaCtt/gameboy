package cartridge

import "github.com/lucactt/gameboy/gameboy/memory"

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
	hCheckStart uint16 = 0x0134
	hCheckEnd   uint16 = 0x014C
)

// Type is used to specify the MBC (if any),
// and if there is additional hardware in the cartridge.
type Type int

// Various cartridge types
const (
	ROM Type = iota
	MBC1
	MBC2
	MBC3
	MBC5
)

// Cartridge represents a GameBoy cartridge.
type Cartridge struct {
	rom *memory.ROM
}

// New creates a new cartridge with the given ROM.
//
// The ROM addresses must start from 0x0000,
// and the total size must at least be the one
// specified by the ROM size memory value,
// otherwise an error will be returned.
//
// Memory space outside the known ROM size won't be touched.
func New(rom *memory.ROM) (*Cartridge, error) {
	return &Cartridge{rom}, nil
}

// Logo returns the value of the Nintendo logo
// contained in the cartridge.
func (c *Cartridge) Logo() [48]byte {
	logo := [48]byte{}

	for i, b := range c.getBytes(logoStart, logoEnd) {
		logo[i] = b
	}

	return logo
}

// Title returns the title of the game in uppercase.
func (c *Cartridge) Title() string {
	return ""
}

// ManufacturerCode returns the manufacturer code in uppercase.
func (c *Cartridge) ManufacturerCode() string {
	return ""
}

// LicenseeCode returns the licensee code in uppercase.
func (c *Cartridge) LicenseeCode() string {
	return ""
}

// DestinationCode returns false if the game is supposed to be sold
// in Japan, true otherwise.
func (c *Cartridge) DestinationCode() bool {
	return false
}

// Type returns the type of the cartridge.
func (c *Cartridge) Type() Type {
	return MBC1
}

func (c *Cartridge) ROMSize() uint16 {
	return 0
}

func (c *Cartridge) RAMSize() uint16 {
	return 0
}

func (c *Cartridge) ROMVersion() int {
	return 0
}

func (c *Cartridge) GlobalChecksum() string {
	return ""
}

func (c *Cartridge) HeaderChecksum() byte {
	return c.getByte(0x014D)
}

// IsValid checks the cartridge validity by
// verifying the correctness of the logo and
// header checksum.
//
// This does not verify the global checksum
// in order to emulate the GameBoy behavior accurately.
func (c *Cartridge) IsValid() bool {
	if c.Logo() != nintendoLogo {
		return false
	}

	var sum byte
	for _, b := range c.getBytes(hCheckStart, hCheckEnd) {
		sum -= b - 1
	}
	return sum == c.HeaderChecksum()
}

func (c *Cartridge) getByte(addr uint16) byte {
	b, err := c.rom.GetByte(addr)

	if err != nil {
		panic(err)
	}

	return b
}

func (c *Cartridge) getBytes(start, end uint16) []byte {
	result := make([]byte, start-end)

	for i := hCheckStart; i <= hCheckEnd; i++ {
		b, err := c.rom.GetByte(i)
		if err != nil {
			panic(err)
		}
		result[i] = b
	}

	return result
}

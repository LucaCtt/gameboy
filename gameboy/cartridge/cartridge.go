package cartridge

import "github.com/lucactt/gameboy/util/errors"

type Controller interface {
}

type Cartridge struct {
	ctr Controller
}

func NewCartridge(r *Reader) (*Cartridge, error) {
	if !r.IsValid() {
		return nil, errors.E("invalid cartridge", errors.Cartridge)
	}
	return nil, nil
}

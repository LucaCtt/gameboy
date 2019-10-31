package memory

import (
	"errors"
	"testing"

	"github.com/lucactt/gameboy/util"
)

type TestSpace struct {
	start    uint16
	end      uint16
	value    byte
	forceErr bool
}

func (s *TestSpace) GetByte(addr uint16) (byte, error) {
	if !s.Accepts(addr) || s.forceErr {
		return 0, errors.New("test")
	}

	return s.value, nil
}

func (s *TestSpace) SetByte(addr uint16, value byte) error {
	if !s.Accepts(addr) || s.forceErr {
		return errors.New("test")
	}

	s.value = value
	return nil
}

func (s *TestSpace) Accepts(addr uint16) bool {
	return addr >= s.start && addr < s.end
}

func TestMemory_GetByte(t *testing.T) {
	tests := []struct {
		name     string
		addr     uint16
		spaceErr bool
		want     byte
		wantErr  bool
	}{
		{"addr in memory", 0x0001, false, 0x11, false},
		{"addr not in memory", 0x1001, false, 0, true},
		{"space error", 0x0001, true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space := &TestSpace{start: 0x0000, end: 0x1000, value: tt.want, forceErr: tt.spaceErr}
			mem := &Memory{}
			mem.AddSpace(space)

			got, err := mem.GetByte(tt.addr)
			util.AssertErr(t, err, tt.wantErr)
			util.AssertEqual(t, got, tt.want)
		})
	}
}

func TestMemory_SetByte(t *testing.T) {
	tests := []struct {
		name     string
		addr     uint16
		spaceErr bool
		want     byte
		wantErr  bool
	}{
		{"addr in memory", 0x0001, false, 0x11, false},
		{"addr not in memory", 0x1001, false, 0, true},
		{"space error", 0x0001, true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space := &TestSpace{start: 0x0000, end: 0x1000, forceErr: tt.spaceErr}
			mem := &Memory{}
			mem.AddSpace(space)

			err := mem.SetByte(tt.addr, 0x11)
			got, err := mem.GetByte(tt.addr)
			util.AssertErr(t, err, tt.wantErr)
			util.AssertEqual(t, got, tt.want)
		})
	}
}

func TestMemory_Accepts(t *testing.T) {
	t.Run("addr in space", func(t *testing.T) {
		space := &TestSpace{start: 0x0000, end: 0x1000}
		mem := &Memory{}
		mem.AddSpace(space)

		got := mem.Accepts(0x0001)
		util.AssertEqual(t, got, true)
	})

	t.Run("addr not in space", func(t *testing.T) {
		space := &TestSpace{start: 0x0000, end: 0x1000}
		mem := &Memory{}
		mem.AddSpace(space)

		got := mem.Accepts(0x1001)
		util.AssertEqual(t, got, false)
	})
}

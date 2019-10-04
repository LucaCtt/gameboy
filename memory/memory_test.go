package memory

import (
	"errors"
	"testing"
)

func assertErr(t *testing.T, got error, want bool) {
	t.Helper()
	if (got != nil) != want {
		t.Errorf("got error %q, want %t", got, want)
	}
}

func assertByte(t *testing.T, got, want byte) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

type TestAddressSpace struct {
	start     uint16
	end       uint16
	getResult byte
	setValue  byte
	wantErr   bool
}

func (s *TestAddressSpace) GetByte(addr uint16) (byte, error) {
	if s.wantErr {
		return 0, errors.New("test")
	}

	if s.setValue != 0 {
		return s.setValue, nil
	}

	return s.getResult, nil
}

func (s *TestAddressSpace) SetByte(addr uint16, value byte) error {
	if s.wantErr {
		return errors.New("test")
	}

	s.setValue = value
	return nil
}

func (s *TestAddressSpace) Accepts(addr uint16) bool {
	return addr >= s.start && addr < s.end
}

func TestMemory_GetByte(t *testing.T) {
	t.Run("addr in space", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000, getResult: 0x11}
		mem := &Memory{}
		mem.AddSpace(space)

		got, err := mem.GetByte(0x0001)
		assertErr(t, err, false)
		assertByte(t, got, 0x11)
	})

	t.Run("addr not in space", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000}
		mem := &Memory{}
		mem.AddSpace(space)

		got, err := mem.GetByte(0x1001)
		assertErr(t, err, true)
		assertByte(t, got, 0)
	})

	t.Run("address space error", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000, wantErr: true}
		mem := &Memory{}
		mem.AddSpace(space)

		got, err := mem.GetByte(0x0001)
		assertErr(t, err, true)
		assertByte(t, got, 0)
	})
}

func TestMemory_SetByte(t *testing.T) {
	t.Run("addr in space", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000}
		mem := &Memory{}
		mem.AddSpace(space)

		err := mem.SetByte(0x0001, 0x11)
		got, err := mem.GetByte(0x0001)
		assertErr(t, err, false)
		assertByte(t, got, 0x11)
	})

	t.Run("addr not in space", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000}
		mem := &Memory{}
		mem.AddSpace(space)

		err := mem.SetByte(0x1001, 0x11)
		assertErr(t, err, true)
	})

	t.Run("address space error", func(t *testing.T) {
		space := &TestAddressSpace{start: 0x0000, end: 0x1000, wantErr: true}
		mem := &Memory{}
		mem.AddSpace(space)

		err := mem.SetByte(0x0001, 0x11)
		assertErr(t, err, true)
	})
}

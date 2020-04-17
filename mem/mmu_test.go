package mem

import (
	"errors"
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

// TestMem is a simple memory that stores a single byte value,
// and can be forced to return an error.
type TestMem struct {
	len      uint16
	value    byte
	forceErr bool
}

func (m *TestMem) GetByte(addr uint16) (byte, error) {
	if !m.Accepts(addr) || m.forceErr {
		return 0, errors.New("test")
	}

	return m.value, nil
}

func (m *TestMem) SetByte(addr uint16, value byte) error {
	if !m.Accepts(addr) || m.forceErr {
		return errors.New("test")
	}

	m.value = value
	return nil
}

func (m *TestMem) Accepts(addr uint16) bool {
	return addr < m.len
}

func TestMem_GetByte(t *testing.T) {
	tests := []struct {
		name      string
		addr      uint16
		spaceErr  bool
		want      byte
		wantErr   bool
		wantPanic bool
	}{
		{"addr in memory", 0x0001, false, 0x11, false, false},
		{"addr not in memory", 0x1000, false, 0, true, false},
		{"space error", 0x0001, true, 0, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("did not want panic")
				}
			}()

			mem := &TestMem{len: 0x1000, value: tt.want, forceErr: tt.spaceErr}
			mmu := &MMU{}
			mmu.AddMem(0x0000, mem)

			got, err := mmu.GetByte(tt.addr)
			assert.Err(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMem_SetByte(t *testing.T) {
	tests := []struct {
		name      string
		addr      uint16
		spaceErr  bool
		want      byte
		wantErr   bool
		wantPanic bool
	}{
		{"addr in memory", 0x0001, false, 0x11, false, false},
		{"addr not in memory", 0x1000, false, 0, true, false},
		{"space error", 0x0001, true, 0, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("did not want panic")
				}
			}()

			mem := &TestMem{len: 0x1000, forceErr: tt.spaceErr}
			mmu := &MMU{}
			mmu.AddMem(0x0000, mem)

			err := mmu.SetByte(tt.addr, 0x11)
			got, err := mmu.GetByte(tt.addr)
			assert.Err(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMem_Accepts(t *testing.T) {
	t.Run("addr in space", func(t *testing.T) {
		mem := &TestMem{len: 0x1000}
		mmu := &MMU{}
		mmu.AddMem(0x0000, mem)

		got := mmu.Accepts(0x0001)
		assert.Equal(t, got, true)
	})

	t.Run("addr not in space", func(t *testing.T) {
		mem := &TestMem{len: 0x1000}
		mmu := &MMU{}
		mmu.AddMem(0x0000, mem)

		got := mmu.Accepts(0x1000)
		assert.Equal(t, got, false)
	})
}

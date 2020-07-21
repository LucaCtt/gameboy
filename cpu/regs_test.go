package cpu

import (
	"testing"

	"github.com/lucactt/gameboy/util/assert"
)

func TestReg_Lo(t *testing.T) {
	r := reg{r: 0x0001}
	got := r.Lo()

	assert.Equal(t, got, byte(0x01))
}

func TestReg_Hi(t *testing.T) {
	r := reg{r: 0x0100}
	got := r.Hi()

	assert.Equal(t, got, byte(0x01))
}

func TestReg_HiLo(t *testing.T) {
	r := reg{r: 0x0100}
	got := r.HiLo()

	assert.Equal(t, got, uint16(0x0100))
}

func TestReg_SetLo(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		r := reg{r: 0x0000, mask: 0xFFF0}
		r.SetLo(byte(0x11))
		got := r.Lo()

		assert.Equal(t, got, byte(0x10))
	})
	t.Run("normal", func(t *testing.T) {
		r := reg{r: 0x0000}
		r.SetLo(byte(0x01))
		got := r.Lo()

		assert.Equal(t, got, byte(0x01))
	})
}

func TestReg_SetHi(t *testing.T) {
	r := reg{r: 0x0000}
	r.SetHi(byte(0x01))
	got := r.Hi()

	assert.Equal(t, got, byte(0x01))
}

func TestReg_Set(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		r := reg{r: 0x0000, mask: 0xFFF0}
		r.Set(0x0111)
		got := r.HiLo()

		assert.Equal(t, got, uint16(0x0110))
	})
	t.Run("normal", func(t *testing.T) {
		r := reg{r: 0x0000}
		r.Set(0x0101)
		got := r.HiLo()

		assert.Equal(t, got, uint16(0x0101))
	})
}

func TestRegs_Z(t *testing.T) {
	tests := []struct {
		name    string
		initial bool
		want    bool
	}{
		{"z initial unset, set", false, false},
		{"z initial unset, unset", false, false},
		{"z initial set, set", true, true},
		{"z initial set, unset", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initial := 0
			if tt.initial {
				initial = 1
			}

			regs := &Regs{AF: reg{r: uint16(initial << 7)}}

			regs.SetZ(tt.want)
			got := regs.Z()

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestRegs_N(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"n is unset", false},
		{"n is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regs := &Regs{}

			regs.SetN(tt.want)
			got := regs.N()

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCPU_H(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"h is unset", false},
		{"h is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regs := &Regs{}

			regs.SetH(tt.want)
			got := regs.H()

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestRegs_C(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"c is unset", false},
		{"c is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regs := &Regs{}

			regs.SetC(tt.want)
			got := regs.C()

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNewRegs(t *testing.T) {
	t.Run("correct initial values", func(t *testing.T) {
		regs := NewRegs()

		assert.Equal(t, regs.AF.HiLo(), defaultAF)
		assert.Equal(t, regs.BC.HiLo(), defaultBC)
		assert.Equal(t, regs.DE.HiLo(), defaultDE)
		assert.Equal(t, regs.HL.HiLo(), defaultHL)
		assert.Equal(t, regs.SP.HiLo(), defaultSP)
		assert.Equal(t, regs.PC.HiLo(), defaultPC)
	})

	t.Run("Z", func(t *testing.T) {
		regs := NewRegs()

		regs.SetZ(true)
		assert.Equal(t, regs.Z(), true)

		regs.SetZ(false)
		assert.Equal(t, regs.Z(), false)
	})

	t.Run("N", func(t *testing.T) {
		regs := NewRegs()

		regs.SetN(true)
		assert.Equal(t, regs.N(), true)

		regs.SetN(false)
		assert.Equal(t, regs.N(), false)
	})

	t.Run("H", func(t *testing.T) {
		regs := NewRegs()

		regs.SetH(true)
		assert.Equal(t, regs.H(), true)

		regs.SetH(false)
		assert.Equal(t, regs.H(), false)
	})

	t.Run("C", func(t *testing.T) {
		regs := NewRegs()

		regs.SetC(true)
		assert.Equal(t, regs.C(), true)

		regs.SetC(false)
		assert.Equal(t, regs.C(), false)
	})
}

package cpu

import (
	"testing"

	"github.com/lucactt/gameboy/util"
)

func TestRegister_Lo(t *testing.T) {
	r := register{r: 0x0001}
	got := r.Lo()

	util.AssertEqual(t, got, byte(0x01))
}

func TestRegister_Hi(t *testing.T) {
	r := register{r: 0x0100}
	got := r.Hi()

	util.AssertEqual(t, got, byte(0x01))
}

func TestRegister_HiLo(t *testing.T) {
	r := register{r: 0x0100}
	got := r.HiLo()

	util.AssertEqual(t, got, uint16(0x0100))
}

func TestRegister_SetLo(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		r := register{r: 0x0000, mask: 0xFFF0}
		r.SetLo(byte(0x11))
		got := r.Lo()

		util.AssertEqual(t, got, byte(0x10))
	})
	t.Run("normal", func(t *testing.T) {
		r := register{r: 0x0000}
		r.SetLo(byte(0x01))
		got := r.Lo()

		util.AssertEqual(t, got, byte(0x01))
	})
}

func TestRegister_SetHi(t *testing.T) {
	r := register{r: 0x0000}
	r.SetHi(byte(0x01))
	got := r.Hi()

	util.AssertEqual(t, got, byte(0x01))
}

func TestRegister_Set(t *testing.T) {
	t.Run("immutable", func(t *testing.T) {
		r := register{r: 0x0000, mask: 0xFFF0}
		r.Set(0x0111)
		got := r.HiLo()

		util.AssertEqual(t, got, uint16(0x0110))
	})
	t.Run("normal", func(t *testing.T) {
		r := register{r: 0x0000}
		r.Set(0x0101)
		got := r.HiLo()

		util.AssertEqual(t, got, uint16(0x0101))
	})
}

func TestCPU_Init(t *testing.T) {
	cpu := New()

	util.AssertEqual(t, cpu.AF.HiLo(), defaultAF)
	util.AssertEqual(t, cpu.BC.HiLo(), defaultBC)
	util.AssertEqual(t, cpu.DE.HiLo(), defaultDE)
	util.AssertEqual(t, cpu.HL.HiLo(), defaultHL)
	util.AssertEqual(t, cpu.SP.HiLo(), defaultSP)
	util.AssertEqual(t, cpu.PC.HiLo(), defaultPC)
}

func TestCPU_Z(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"z is unset", false},
		{"z is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{}

			cpu.SetZ(tt.want)
			got := cpu.Z()

			util.AssertEqual(t, got, tt.want)
		})
	}
}

func TestCPU_N(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"n is unset", false},
		{"n is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{}

			cpu.SetN(tt.want)
			got := cpu.N()

			util.AssertEqual(t, got, tt.want)
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
			cpu := &CPU{}

			cpu.SetH(tt.want)
			got := cpu.H()

			util.AssertEqual(t, got, tt.want)
		})
	}
}

func TestCPU_C(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"c is unset", false},
		{"c is set", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{}

			cpu.SetC(tt.want)
			got := cpu.C()

			util.AssertEqual(t, got, tt.want)
		})
	}
}

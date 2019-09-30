package cpu

import "testing"

func assertBytesEqual(t *testing.T, got, want byte) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertBitsEqual(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestRegister_Lo(t *testing.T) {
	r := register{0x0001}
	got := r.Lo()

	assertBytesEqual(t, got, byte(0x01))
}

func TestRegister_Hi(t *testing.T) {
	r := register{0x0100}
	got := r.Hi()

	assertBytesEqual(t, got, byte(0x01))
}

func TestRegister_HiLo(t *testing.T) {
	r := register{0x0100}
	got := r.HiLo()

	assertBytesEqual(t, byte(got), byte(0x00))
	assertBytesEqual(t, byte(got>>8), byte(0x01))
}

func TestRegister_SetLo(t *testing.T) {
	r := register{0x0000}
	r.SetLo(byte(0x01))
	got := r.Lo()

	assertBytesEqual(t, got, byte(0x01))
}

func TestRegister_SetHi(t *testing.T) {
	r := register{0x0000}
	r.SetHi(byte(0x01))
	got := r.Hi()

	assertBytesEqual(t, got, byte(0x01))
}

func TestRegister_Set(t *testing.T) {
	r := register{0x0000}
	r.Set(0x0101)
	got := r.HiLo()

	assertBytesEqual(t, byte(got), byte(0x01))
	assertBytesEqual(t, byte(got>>8), byte(0x01))
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

			assertBitsEqual(t, got, tt.want)
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

			assertBitsEqual(t, got, tt.want)
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

			assertBitsEqual(t, got, tt.want)
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

			assertBitsEqual(t, got, tt.want)
		})
	}
}

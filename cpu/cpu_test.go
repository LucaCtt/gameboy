package cpu

import "testing"

func assertBitsEqual(t *testing.T, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
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

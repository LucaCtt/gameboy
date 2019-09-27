package cpu

import "testing"

func assertBitsEqual(t *testing.T, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestCPU_Z(t *testing.T) {
	t.Run("z is unset", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetZ(false)
		got := cpu.Z()

		assertBitsEqual(t, got, false)
	})

	t.Run("z is set", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetZ(true)
		got := cpu.Z()

		assertBitsEqual(t, got, true)
	})
}

func TestCPU_N(t *testing.T) {
	t.Run("n is unset", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetN(false)
		got := cpu.N()

		assertBitsEqual(t, got, false)
	})

	t.Run("n is set", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetN(true)
		got := cpu.N()

		assertBitsEqual(t, got, true)
	})
}

func TestCPU_H(t *testing.T) {
	t.Run("h is unset", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetH(false)
		got := cpu.H()

		assertBitsEqual(t, got, false)
	})

	t.Run("h is set", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetH(true)
		got := cpu.H()

		assertBitsEqual(t, got, true)
	})
}

func TestCPU_C(t *testing.T) {
	t.Run("c is unset", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetC(false)
		got := cpu.C()

		assertBitsEqual(t, got, false)
	})

	t.Run("c is set", func(t *testing.T) {
		cpu := &CPU{}

		cpu.SetC(true)
		got := cpu.C()

		assertBitsEqual(t, got, true)
	})
}

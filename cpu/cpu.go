// Package cpu implements a complete GameBoy CPU.
package cpu

import (
	"github.com/lucactt/gameboy/mem"
	"github.com/lucactt/gameboy/util/errors"
)

// CPU represents a GameBoy CPU.
type CPU struct {
	Mem      mem.Mem
	Regs     *Regs
	StateMgr *StateMgr
	InstrSet *InstrSet
}

// New creates a new CPU.
func New(mem mem.Mem) *CPU {
	regs := NewRegs()
	stateMgr := NewStateMgr()
	instrSet := NewInstrSet(regs, mem, stateMgr)

	return &CPU{mem, regs, stateMgr, instrSet}
}

// Tick runs the instruction found in the memory at the address contained in PC,
// and returns the number of clock cycles used by that instruction.
func (c *CPU) Tick() (int, error) {
	pc := c.Regs.PC.HiLo()
	opCode, err := c.Mem.GetByte(pc)
	if err != nil {
		return 0, errors.E("get opcode failed", err, errors.CPU)
	}

	len, cycles := c.InstrSet.NoPrefix[opCode]()
	c.Regs.PC.Set(pc + uint16(len))

	return cycles, nil
}

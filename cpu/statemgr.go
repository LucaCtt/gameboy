package cpu

// State represents a CPU state.
type State string

// CPU states.
const (
	Running State = "running"
	Stopped State = "stopped"
	Halted  State = "halted"
)

// StateMgr manages the states of the CPU.
type StateMgr struct {
	current State
	ime     bool // Interrupt Master Enable
}

// NewStateMgr creates a new StateMgr.
func NewStateMgr() *StateMgr {
	return &StateMgr{current: Running, ime: true}
}

// State returns the current CPU state.
func (s *StateMgr) State(n State) State {
	return s.current
}

// SetState sets the current CPU state.
func (s *StateMgr) SetState(n State) {
	s.current = n
}

// InterruptsEnabled returns the current IME state.
func (s *StateMgr) InterruptsEnabled(n State) bool {
	return s.ime
}

// SetIME enables or disables interrupts handling.
func (s *StateMgr) SetIME(v bool) {
	s.ime = v
}

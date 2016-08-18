package user

const (
	FSM_PLAYING = iota
	FSM_STOP
	FSM_PAUSE
)

//FSM A finite-state machine
type FSM struct {
	State byte
}

// Play media
func (fsm *FSM) Play() {}

// Pause media
func (fsm *FSM) Pause() {}

// Stop media
func (fsm *FSM) Stop() {}

func NewFSM() *FSM {
	return &FSM{State: FSM_STOP}
}

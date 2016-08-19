package user

import "github.com/doublemo/koala/media"

const (
	FSM_PLAYING = iota
	FSM_STOP
	FSM_PAUSE
)

//FSM A finite-state machine
type FSM struct {
	State        byte
	Parameters   *media.StreamParameters
	MediaName    string
	TaskID       string
	MediaSession *media.MediaSession
}

// Play media
func (fsm *FSM) Play(ssrc uint32) {
	if fsm.State == FSM_PLAYING {
		return
	}

	fsm.MediaSession.Play(ssrc, fsm.Parameters)

	fsm.State = FSM_PLAYING
}

// Pause media
func (fsm *FSM) Pause() {}

// Stop media
func (fsm *FSM) Stop() {}

func NewFSM(mediaName, taskID string, parameters *media.StreamParameters, msess *media.MediaSession) *FSM {
	return &FSM{
		State:        FSM_STOP,
		Parameters:   parameters,
		MediaName:    mediaName,
		TaskID:       taskID,
		MediaSession: msess,
	}
}

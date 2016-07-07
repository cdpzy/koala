package mediasession

import (
    "fmt"
)

type ServerMediaSubSession struct {
    Id               int
    FileName         string
    FileSize         int64
}


func (serverMediaSubSession *ServerMediaSubSession) IncrementTrackId() {
    serverMediaSubSession.Id ++
}

func (serverMediaSubSession *ServerMediaSubSession) GetTrackId() string {
    return fmt.Sprintf("track%d", serverMediaSubSession.Id)
}
package helper

import (
    "time"
)

type Time struct {
    SEC   int64
    USEC  int64
}

func GetNowTime() *Time {
    nsec  := time.Now().UnixNano()
    v     := &Time{}
    v.SEC  = nsec / 1000000000
    v.USEC = nsec % (v.SEC * 1000000000)
    return v
}
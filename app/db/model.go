package db

import (
  "time"
)

type TimerInterface interface {
  IsEnded() bool
  Start(isWork bool) time.Time
}

type Timer struct {
  Id uint64
  EndTimeUnixSec int64 `json:"EndTime" db:"EndTime"`
  IsBreak bool `json:"IsBreak" db:"IsBreak"`
}

func (t Timer) IsEnded() bool {
  now := time.Now()
  return now.Unix() > t.EndTimeUnixSec
}

func (t Timer) Start() Timer {
  var nextEnd time.Time
  var isBreak bool
  if t.IsBreak {
    nextEnd = time.Now().Add(time.Minute * 25)
    isBreak = false
  } else {
    nextEnd = time.Now().Add(time.Minute * 5)
    isBreak = true
  }
  t.EndTimeUnixSec = nextEnd.Unix()
  t.IsBreak = isBreak
  return t
}

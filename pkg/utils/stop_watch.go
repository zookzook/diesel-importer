package utils

import (
	"fmt"
	"time"
)

type StopWatch struct {
	Message  string
	T        time.Time
	Duration time.Duration
}

func NewStopWatch() StopWatch {
	return StopWatch{}
}

func (s *StopWatch) Start(msg string) {
	s.T = time.Now()
	s.Message = msg
}

func (s *StopWatch) Stop() {
	s.Duration = time.Since(s.T)
}

func (s *StopWatch) String() string {
	return fmt.Sprintf("%s took %d", s.Message, s.Duration)
}

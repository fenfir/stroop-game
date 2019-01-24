package main

import "time"

type SecondsTimer struct {
	Timer *time.Timer
	End   time.Time
}

func NewSecondsTimer(seconds time.Duration) *SecondsTimer {
	return &SecondsTimer{
		Timer: time.NewTimer(seconds * time.Second),
		End:   time.Now().Add(seconds * time.Second),
	}
}

func (s *SecondsTimer) Stop() {
	s.Timer.Stop()
}

func (s *SecondsTimer) Remaining() time.Duration {
	return s.End.Sub(time.Now())
}

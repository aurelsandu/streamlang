package clock

import "time"

// Clock permite injecție de timp pentru operatori (pt. teste deterministe).
type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
}

type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

package clock

import "time"

type realClock struct{}

func NewReal() Clock { return realClock{} }

func (realClock) Now() time.Time { return time.Now() }

func (realClock) NewTimer(d time.Duration) Timer {
	return &timerWrap{t: time.NewTimer(d)}
}
func (realClock) NewTicker(d time.Duration) Ticker {
	return &tickerWrap{t: time.NewTicker(d)}
}

type timerWrap struct{ t *time.Timer }
func (w *timerWrap) C() <-chan time.Time          { return w.t.C }
func (w *timerWrap) Stop() bool                   { return w.t.Stop() }
func (w *timerWrap) Reset(d time.Duration) bool   { return w.t.Reset(d) }

type tickerWrap struct{ t *time.Ticker }
func (w *tickerWrap) C() <-chan time.Time { return w.t.C }
func (w *tickerWrap) Stop()               { w.t.Stop() }

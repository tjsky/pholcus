package spider

import (
	"sync"
	"time"

	"github.com/andeya/pholcus/logs"
)

// Timer manages a collection of named clocks (countdown timers or alarms).
type Timer struct {
	setting map[string]*Clock
	closed  bool
	sync.RWMutex
}

func newTimer() *Timer {
	return &Timer{
		setting: make(map[string]*Clock),
	}
}

// sleep blocks until the named timer fires and reports whether it can still be used.
func (self *Timer) sleep(id string) bool {
	self.RLock()
	if self.closed {
		self.RUnlock()
		return false
	}

	c, ok := self.setting[id]
	self.RUnlock()
	if !ok {
		return false
	}

	c.sleep()

	self.RLock()
	defer self.RUnlock()
	if self.closed {
		return false
	}
	_, ok = self.setting[id]

	return ok
}

// set configures a timer. When bell is nil, tol is a countdown sleep duration;
// otherwise tol specifies the wake-up occurrence (the tol-th bell from now).
func (self *Timer) set(id string, tol time.Duration, bell *Bell) bool {
	self.Lock()
	defer self.Unlock()
	if self.closed {
		logs.Log.Critical("*** timer [%s]: failed to set, timer system is closed ***", id)
		return false
	}
	c, ok := newClock(id, tol, bell)
	if !ok {
		logs.Log.Critical("*** timer [%s]: failed to set, invalid parameters ***", id)
		return ok
	}
	self.setting[id] = c
	logs.Log.Critical("*** timer [%s]: set successfully ***", id)
	return ok
}

// drop cancels all timers and marks the Timer as closed.
func (self *Timer) drop() {
	self.Lock()
	defer self.Unlock()
	self.closed = true
	for _, c := range self.setting {
		c.wake()
	}
	self.setting = make(map[string]*Clock)
}

type (
	// Clock represents a single alarm or countdown timer.
	Clock struct {
		id    string
		typ   int           // mode: A (alarm) or T (countdown)
		tol   time.Duration // countdown duration, or alarm occurrence count
		bell  *Bell         // alarm time-of-day (nil for countdown mode)
		timer *time.Timer
	}
	// Bell specifies a time-of-day for alarm mode.
	Bell struct {
		Hour int
		Min  int
		Sec  int
	}
)

const (
	A = iota // alarm mode
	T        // countdown mode
)

// newClock creates a Clock. When bell is nil, tol is a countdown duration;
// otherwise tol specifies the wake-up occurrence.
func newClock(id string, tol time.Duration, bell *Bell) (*Clock, bool) {
	if tol <= 0 {
		return nil, false
	}
	if bell == nil {
		return &Clock{
			id:    id,
			typ:   T,
			tol:   tol,
			timer: newT(),
		}, true
	}
	if !(bell.Hour >= 0 && bell.Hour < 24 && bell.Min >= 0 && bell.Min < 60 && bell.Sec >= 0 && bell.Sec < 60) {
		return nil, false
	}
	return &Clock{
		id:    id,
		typ:   A,
		tol:   tol,
		bell:  bell,
		timer: newT(),
	}, true
}

func (self *Clock) sleep() {
	d := self.duration()
	self.timer.Reset(d)
	t0 := time.Now()
	logs.Log.Critical("*** timer <%s> sleeping %v, scheduled wake at %v ***", self.id, d, t0.Add(d).Format("2006-01-02 15:04:05"))
	<-self.timer.C
	t1 := time.Now()
	logs.Log.Critical("*** timer <%s> woke at %v, actual sleep %v ***", self.id, t1.Format("2006-01-02 15:04:05"), t1.Sub(t0))
}

func (self *Clock) wake() {
	self.timer.Reset(0)
}

func (self *Clock) duration() time.Duration {
	switch self.typ {
	case A:
		t := time.Now()
		year, month, day := t.Date()
		bell := time.Date(year, month, day, self.bell.Hour, self.bell.Min, self.bell.Sec, 0, time.Local)
		if bell.Before(t) {
			bell = bell.Add(time.Hour * 24 * self.tol)
		} else {
			bell = bell.Add(time.Hour * 24 * (self.tol - 1))
		}
		return bell.Sub(t)
	case T:
		return self.tol
	}
	return 0
}

func newT() *time.Timer {
	t := time.NewTimer(0)
	<-t.C
	return t
}

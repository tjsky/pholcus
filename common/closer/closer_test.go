package closer

import (
	"errors"
	"testing"
)

type mockCloser struct {
	err error
}

func (m *mockCloser) Close() error { return m.err }

func TestLogClose_NoError(t *testing.T) {
	var logged bool
	LogClose(&mockCloser{}, func(string, ...interface{}) { logged = true })
	if logged {
		t.Error("expected no log call when Close succeeds")
	}
}

func TestLogClose_WithError(t *testing.T) {
	var msg string
	LogClose(&mockCloser{err: errors.New("boom")}, func(format string, args ...interface{}) {
		msg = format
		_ = args
	})
	if msg == "" {
		t.Error("expected log call when Close fails")
	}
}

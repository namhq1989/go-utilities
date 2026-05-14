package logger

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

type captureHook struct {
	entries []*logrus.Entry
}

func (h *captureHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *captureHook) Fire(entry *logrus.Entry) error {
	h.entries = append(h.entries, entry)
	return nil
}

// newTestLogger builds a Logger wired to a fresh logrus.Logger so tests are
// isolated from the package-level singleton in logger.go.
func newTestLogger(t *testing.T) (*Logger, *captureHook, *bytes.Buffer) {
	t.Helper()
	base := logrus.New()
	buf := &bytes.Buffer{}
	base.SetOutput(buf)
	base.SetLevel(logrus.DebugLevel)
	hook := &captureHook{}
	base.AddHook(hook)
	return &Logger{logger: base.WithFields(Fields{"source": "test"})}, hook, buf
}

func TestError_AttachesErrorViaWithError(t *testing.T) {
	l, hook, _ := newTestLogger(t)
	wantErr := errors.New("boom")

	l.Error("failed to do thing", wantErr, Fields{"id": "x"})

	if len(hook.entries) != 1 {
		t.Fatalf("hook entries = %d, want 1", len(hook.entries))
	}
	entry := hook.entries[0]
	if entry.Message != "failed to do thing" {
		t.Errorf("Message = %q, want %q", entry.Message, "failed to do thing")
	}
	got, ok := entry.Data[logrus.ErrorKey].(error)
	if !ok {
		t.Fatalf("logrus.ErrorKey not set or not an error: %T %v", entry.Data[logrus.ErrorKey], entry.Data[logrus.ErrorKey])
	}
	if !errors.Is(got, wantErr) {
		t.Errorf("attached error %v does not match %v", got, wantErr)
	}
	if entry.Data["id"] != "x" {
		t.Errorf("user field id missing: %v", entry.Data)
	}
}

func TestError_NilErrorDoesNotSetErrorKey(t *testing.T) {
	l, hook, _ := newTestLogger(t)
	l.Error("just a message", nil, nil)
	if len(hook.entries) != 1 {
		t.Fatalf("hook entries = %d", len(hook.entries))
	}
	if _, ok := hook.entries[0].Data[logrus.ErrorKey]; ok {
		t.Errorf("logrus.ErrorKey should not be set when err is nil")
	}
}

func TestError_RedactsFields(t *testing.T) {
	l, hook, _ := newTestLogger(t)
	l.Error("failed", errors.New("e"), Fields{
		"password": "p",
		"email":    "alice@example.com",
		"user_id":  "u_123",
	})

	entry := hook.entries[0]
	if entry.Data["password"] != RedactedPlaceholder {
		t.Errorf("password not redacted: %v", entry.Data["password"])
	}
	if entry.Data["email"] != "alic*********.com" {
		t.Errorf("email not masked: %v", entry.Data["email"])
	}
	if entry.Data["user_id"] != "u_123" {
		t.Errorf("user_id mutated: %v", entry.Data["user_id"])
	}
}

func TestSilentError_SetsMarker(t *testing.T) {
	l, hook, _ := newTestLogger(t)
	l.SilentError("failed silently", errors.New("e"), Fields{"k": "v"})

	if len(hook.entries) != 1 {
		t.Fatalf("hook entries = %d", len(hook.entries))
	}
	entry := hook.entries[0]
	silent, ok := entry.Data[MarkerSilent].(bool)
	if !ok || !silent {
		t.Errorf("MarkerSilent not set to true: %v", entry.Data[MarkerSilent])
	}
	if _, ok := entry.Data[logrus.ErrorKey].(error); !ok {
		t.Errorf("error still expected to be attached")
	}
}

func TestSilentError_RedactsFields(t *testing.T) {
	l, hook, _ := newTestLogger(t)
	l.SilentError("oops", nil, Fields{"authorization": "Bearer x"})

	if hook.entries[0].Data["authorization"] != RedactedPlaceholder {
		t.Errorf("authorization not redacted: %v", hook.entries[0].Data["authorization"])
	}
}

func TestInfoWarnDebug_RedactFields(t *testing.T) {
	l, _, _ := newTestLogger(t)
	hook := &captureHook{}
	// capture across all levels
	base := logrus.New()
	base.SetLevel(logrus.DebugLevel)
	base.SetOutput(&bytes.Buffer{})
	allLevels := &allLevelHook{}
	base.AddHook(allLevels)
	l = &Logger{logger: base.WithFields(Fields{})}

	l.Info("i", Fields{"token": "t"})
	l.Warn("w", Fields{"password": "p"})
	l.Debug("d", Fields{"secret": "s"})

	if len(allLevels.entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(allLevels.entries))
	}
	if allLevels.entries[0].Data["token"] != RedactedPlaceholder {
		t.Errorf("Info: token not redacted")
	}
	if allLevels.entries[1].Data["password"] != RedactedPlaceholder {
		t.Errorf("Warn: password not redacted")
	}
	if allLevels.entries[2].Data["secret"] != RedactedPlaceholder {
		t.Errorf("Debug: secret not redacted")
	}
	_ = hook
}

type allLevelHook struct{ entries []*logrus.Entry }

func (h *allLevelHook) Levels() []logrus.Level { return logrus.AllLevels }
func (h *allLevelHook) Fire(e *logrus.Entry) error {
	h.entries = append(h.entries, e)
	return nil
}

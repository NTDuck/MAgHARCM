package pattern

import (
	"errors"
	"testing"
)

func TestMustSuccess(t *testing.T) {
	val := Must("hello", nil)
	if val != "hello" {
		t.Fatalf("expected hello, got %s", val)
	}
}

func TestMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on non-nil error")
		}
	}()
	_ = Must("", errors.New("test failure"))
}

func TestMust0Success(t *testing.T) {
	Must0(nil)
}

func TestMust0Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on non-nil error")
		}
	}()
	Must0(errors.New("test failure"))
}

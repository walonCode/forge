package utils

import "testing"

func TestGetEnvString(t *testing.T) {
	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("TEST_STR", "hello")
		if got := GetEnvString("TEST_STR", "fallback"); got != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
	})

	t.Run("returns default when unset", func(t *testing.T) {
		if got := GetEnvString("TEST_STR_MISSING", "fallback"); got != "fallback" {
			t.Fatalf("got %q, want %q", got, "fallback")
		}
	})

	t.Run("returns default when blank", func(t *testing.T) {
		t.Setenv("TEST_STR_BLANK", "   ")
		if got := GetEnvString("TEST_STR_BLANK", "fallback"); got != "fallback" {
			t.Fatalf("got %q, want %q", got, "fallback")
		}
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Run("parses a valid int", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		if got := GetEnvInt("TEST_INT", 7); got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})

	t.Run("returns default on invalid int", func(t *testing.T) {
		t.Setenv("TEST_INT_BAD", "not-a-number")
		if got := GetEnvInt("TEST_INT_BAD", 7); got != 7 {
			t.Fatalf("got %d, want 7", got)
		}
	})

	t.Run("returns default when unset", func(t *testing.T) {
		if got := GetEnvInt("TEST_INT_MISSING", 7); got != 7 {
			t.Fatalf("got %d, want 7", got)
		}
	})
}

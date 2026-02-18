package testutils

import "testing"

func AssertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatalf("wanted an error but didn't get one")
	}
	if got.Error() != want.Error() {
		t.Fatalf("found wrong error: got %q want %q", got, want)
	}
}

func AssertHasError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("wanted an error but didn't get one")
	}
}

func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got error %q but didnt expect one", got)
	}
}

// Deprecated in favor of AssertEqual[T comparable]
func AssertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

// Deprecated in favor of AssertEqual[T comparable]
func AssertInts(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Fatalf("didn't want %v", got)
	}
}

func AssertTrue(t testing.TB, got bool) {
	t.Helper()
	if got != true {
		t.Fatalf("got %t, want true", got)
	}
}

func AssertFalse(t testing.TB, got bool) {
	t.Helper()
	if got != false {
		t.Fatalf("got %t, want false", got)
	}
}

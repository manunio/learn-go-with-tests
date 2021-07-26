package di

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Manu")

	got := buffer.String()
	want := "Hello, Manu"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

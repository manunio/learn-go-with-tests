package main

import "testing"

func TestHello(t *testing.T) {

	t.Run("say hello to people", func(t *testing.T) {
		got := Hello("Manu")
		want := "Hello, Manu"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("say 'Hello, world' when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, world"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

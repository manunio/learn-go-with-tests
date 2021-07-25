package main

import "testing"

func TestHello(t *testing.T) {

	assertCorretMessage := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("in French",func(t *testing.T) {
		got := Hello("Manu", "French")
		want := "Bonjour, Manu"
		assertCorretMessage(t, got, want)
	})

	t.Run("in Spanish",func(t *testing.T) {
		got := Hello("Manu", "Spanish")
		want := "Hola, Manu"
		assertCorretMessage(t, got, want)
	})

	t.Run("say hello to people", func(t *testing.T) {
		got := Hello("Manu","")
		want := "Hello, Manu"

		assertCorretMessage(t, got, want)

	})

	t.Run("say 'Hello, world' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, world"

		assertCorretMessage(t, got, want)
	})

}

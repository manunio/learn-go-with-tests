package integers

import (
	"fmt"
	"testing"
)


func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// output: 6
}

func TestAdder(t *testing.T) {

	assertCorretMessage := func(t testing.TB, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got '%d' want '%d'", got, want)
		}
	}

	t.Run("add 2 numbers", func(t *testing.T) {
		got := Add(2,2)
		want := 4
		assertCorretMessage(t, got, want)
	})

	
}
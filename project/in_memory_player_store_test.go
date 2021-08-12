package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestInMemoryPlayerStoreConcurrently(t *testing.T) {
	t.Run("RecordWin runs safe concurrently", func(t *testing.T) {
		memoryPlayerStore := NewInMemoryPlayerStore()
		wantedStore := 1000

		var wg sync.WaitGroup
		wg.Add(wantedStore)

		for i := 0; i < wantedStore; i++ {
			i := i
			go func() {
				memoryPlayerStore.RecordWin(fmt.Sprint(i))
				wg.Done()
			}()
		}
		wg.Wait()
		gotStore := len(memoryPlayerStore.store)
		if gotStore != wantedStore {
			t.Errorf("got %d want %d ", gotStore, wantedStore)
		}
	})

	// TODO: (manunio) investigate this, as i'm not sure about the authenticity of this test.
	t.Run("GetPlayerScore runs safe concurrently", func(t *testing.T) {
		memoryPlayerStore := NewInMemoryPlayerStore()
		wantedStore := 1000

		for i := 0; i < wantedStore; i++ {
			memoryPlayerStore.RecordWin(fmt.Sprint(i))
		}

		var wg sync.WaitGroup
		wg.Add(wantedStore)

		for i := 0; i < wantedStore; i++ {
			i := i
			go func() {
				memoryPlayerStore.GetPlayerScore(fmt.Sprint(i))
				wg.Done()
			}()
		}

		wg.Wait()
		gotStore := len(memoryPlayerStore.store)
		if gotStore != wantedStore {
			t.Errorf("got %d want %d ", gotStore, wantedStore)
		}
	})
}

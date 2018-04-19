package cmd

import (
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	group := sync.WaitGroup{}
	group.Add(1)

	go func() {
		out, err := Run("sh", "-c", "env && sleep 10 && echo hello")
		t.Log(out, err)
	}()

	go func() {
		for i := 0; i < 100; i++ {
			t.Log("hello")
		}
	}()

	group.Wait()
}

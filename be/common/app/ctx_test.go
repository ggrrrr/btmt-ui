package app

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func work(ctx context.Context, data string) {
	while := true
	ss := time.Duration(len(data))
	for while {
		fmt.Printf("%s start...\n", data)
		select {
		case <-ctx.Done():
			fmt.Printf("%s ctx cancel\n", data)
			while = false

		case <-time.After(ss * time.Second):
			fmt.Printf("%s timer done.\n", data)

		}

		fmt.Printf("%s select done.\n", data)
	}
	fmt.Printf("%s done loop.\n", data)
}

func TestMain(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	batch := []string{"1", "22", "big work"}
	wg := sync.WaitGroup{}
	wg.Add(len(batch))
	go func() {
		<-time.After(2 * time.Second)
		fmt.Println("send kill")
		cancel()
		fmt.Println("canceled.")
	}()
	for i := range batch {
		go func(d string) {
			defer wg.Done()
			work(ctx, d)
		}(batch[i])
	}
	wg.Wait()
}

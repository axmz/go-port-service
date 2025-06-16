package graceful

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Operation = func(ctx context.Context) error

func Shutdown(timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", key)
				if err := op(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", key, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", key)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

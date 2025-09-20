package worker

import (
	"context"
	"sync"
	"urlcheck/internal/checker"
	"urlcheck/internal/types"
)

type Worker struct {
	MaxWorkers int
}

func (w *Worker) validateMaxWorkers() {
	if w.MaxWorkers <= 0 {
		w.MaxWorkers = 1
	}
}

func (w *Worker) Run(
	ctx context.Context,
	c checker.Checker,
	urls []string,
	callback func(current, total int, result *types.Result),
) error {
	w.validateMaxWorkers()

	results := make(chan *types.Result, len(urls))
	urlChan := make(chan string, len(urls))

	for _, u := range urls {
		urlChan <- u
	}
	close(urlChan)

	var wg sync.WaitGroup
	wg.Add(w.MaxWorkers)

	for i := 0; i < w.MaxWorkers; i++ {
		go func() {
			defer wg.Done()
			for url := range urlChan {
				result := c.Check(url)
				results <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	processedCount := 0
	total := len(urls)

	for {
		select {
		case res, ok := <-results:
			if !ok {
				return nil
			}
			processedCount++
			callback(processedCount, total, res)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

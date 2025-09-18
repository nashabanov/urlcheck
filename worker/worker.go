package worker

import (
	"context"
	"sync"
	"urlcheck/checker"
)

type Worker struct {
	MaxWorkers int
}

func (w *Worker) Run(ctx context.Context, c checker.Checker, urls []string) []*checker.Result {
	if w.MaxWorkers <= 0 {
		w.MaxWorkers = 1
	}

	results := make(chan *checker.Result, len(urls))
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

	var allResults []*checker.Result
	for {
		select {
		case res, ok := <-results:
			if !ok {
				return allResults
			}
			allResults = append(allResults, res)
		case <-ctx.Done():
			return allResults
		}
	}
}

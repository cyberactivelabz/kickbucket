package engine

import (
	"context"
	"errors"
	"runtime"
	"sync"
)

// Scan checks candidates concurrently and returns one result per successful check.
func Scan(ctx context.Context, provider Provider, candidates []string, workers int, limiter *Limiter) ([]Result, error) {
	if provider == nil {
		return nil, errors.New("provider is required")
	}
	if workers < 1 {
		workers = runtime.NumCPU()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type checked struct {
		name   string
		access Access
	}
	jobs := make(chan string)
	results := make(chan checked)
	var workerGroup sync.WaitGroup
	for i := 0; i < workers; i++ {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case name, ok := <-jobs:
					if !ok {
						return
					}
					if err := limiter.Wait(ctx); err != nil {
						return
					}
					access, err := provider.Check(ctx, name)
					if err == nil {
						select {
						case results <- checked{name: name, access: access}:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, candidate := range candidates {
			select {
			case jobs <- candidate:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		workerGroup.Wait()
		close(results)
	}()

	output := make([]Result, 0, len(candidates))
	for result := range results {
		output = append(output, Result{Provider: provider.Name(), Name: result.name, Access: result.access})
	}
	return output, ctx.Err()
}

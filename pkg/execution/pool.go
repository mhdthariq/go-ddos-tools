package execution

import (
	"context"
	"sync"
	"time"

	"github.com/go-ddos-tools/pkg/core"
)

// RunAttack executes an attack using a worker pool
func RunAttack(ctx context.Context, attacker core.Attacker, threads int) {
	var wg sync.WaitGroup
	
	// Create a cancelable context if not already done, but usually passed in
	// For this helper, we assume ctx is already managed by the caller (with timeout)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Execute one iteration
					if err := attacker.Attack(ctx); err != nil {
						// On error, we might want to backoff slightly or just continue
						// For high-performance stress testing, usually we just continue
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		}()
	}

	wg.Wait()
}

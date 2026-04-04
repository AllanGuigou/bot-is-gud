package notifier

import (
	"context"
	"time"
)

// Schedule calls f on a periodic schedule with an optional manual trigger.
// p is the period, o is the offset from midnight, trigger allows immediate execution.
func Schedule(ctx context.Context, p time.Duration, o time.Duration, trigger <-chan struct{}, f func(time.Time)) {
	first := time.Now().Truncate(p).Add(o)

	if first.Before(time.Now()) {
		first = first.Add(p)
	}

	c := time.After(time.Until(first))
	// receiving from a nil channel blocks forever
	t := &time.Ticker{C: nil}

	for {
		select {
		case v := <-c:
			t = time.NewTicker(p)
			f(v)
		case v := <-t.C:
			f(v)
		case <-trigger:
			f(time.Now())
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

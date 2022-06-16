package birthday

import (
	"context"
	"time"
)

// calls function `f` with a period `p` offset by `o`
// source: https://stackoverflow.com/a/56156860/11192934
func Schedule(ctx context.Context, p time.Duration, o time.Duration, f func(time.Time)) {
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
			// the ticker has to be started before f as it can take some time to finish
			t = time.NewTicker(p)
			f(v)
		case v := <-t.C:
			f(v)
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

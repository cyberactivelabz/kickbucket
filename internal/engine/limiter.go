package engine

import (
	"context"
	"time"
)

// Limiter spaces requests at the requested rate. A zero rate disables limiting.
type Limiter struct {
	tokens <-chan time.Time
	stop   func()
}

func NewLimiter(perSecond float64) *Limiter {
	if perSecond <= 0 {
		return &Limiter{}
	}
	interval := time.Duration(float64(time.Second) / perSecond)
	ticker := time.NewTicker(interval)
	return &Limiter{tokens: ticker.C, stop: ticker.Stop}
}

// Stop releases the limiter's ticker resources.
func (l *Limiter) Stop() {
	if l != nil && l.stop != nil {
		l.stop()
	}
}

func (l *Limiter) Wait(ctx context.Context) error {
	if l == nil || l.tokens == nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
		return nil
	}
}

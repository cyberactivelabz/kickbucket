package engine

import (
	"context"
	"testing"
)

type testProvider struct{}

func (testProvider) Name() string                                  { return "test" }
func (testProvider) Candidates(string, []string) []string          { return nil }
func (testProvider) Check(context.Context, string) (Access, error) { return AccessPrivate, nil }

func TestScan(t *testing.T) {
	got, err := Scan(context.Background(), testProvider{}, []string{"one", "two", "three"}, 2, nil)
	if err != nil || len(got) != 3 {
		t.Fatalf("Scan() returned %d results and error %v", len(got), err)
	}
}

func TestLimiterHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := NewLimiter(1).Wait(ctx); err == nil {
		t.Fatal("Wait() returned nil for cancelled context")
	}
}

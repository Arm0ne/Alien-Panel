package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPolicyDelayFollowsExponentialSchedule(t *testing.T) {
	policy := Policy{MaxAttempts: 4, InitialBackoff: 30 * time.Second, MaxBackoff: 10 * time.Minute}
	for attempt, want := range map[int]time.Duration{1: 30 * time.Second, 2: 60 * time.Second, 3: 120 * time.Second} {
		if got := policy.Delay(attempt); got != want {
			t.Fatalf("Delay(%d) = %s, want %s", attempt, got, want)
		}
	}
}

func TestPolicySupportsOperationalSchedule(t *testing.T) {
	policy := Policy{MaxAttempts: 4, InitialBackoff: 30 * time.Second, MaxBackoff: 10 * time.Minute, Schedule: []time.Duration{30 * time.Second, 2 * time.Minute, 10 * time.Minute}}
	for attempt, want := range map[int]time.Duration{1: 30 * time.Second, 2: 2 * time.Minute, 3: 10 * time.Minute} {
		if got := policy.Delay(attempt); got != want {
			t.Fatalf("scheduled Delay(%d) = %s, want %s", attempt, got, want)
		}
	}
}

func TestPolicyDelayCapsAtMaximum(t *testing.T) {
	policy := Policy{MaxAttempts: 8, InitialBackoff: 4 * time.Second, MaxBackoff: 5 * time.Second}
	if got := policy.Delay(3); got != 5*time.Second {
		t.Fatalf("Delay(3) = %s, want 5s cap", got)
	}
}

func TestWaitHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Wait(ctx, time.Hour)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait() error = %v, want context.Canceled", err)
	}
}

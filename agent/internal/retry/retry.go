// Package retry contains the small, context-aware retry primitive shared by
// the Agent's remote clients. It deliberately does not decide which errors
// are safe to retry; callers know whether a failed request is idempotent and
// can classify HTTP responses accordingly.
package retry

import (
	"context"
	"time"
)

// Policy controls the maximum number of attempts and the exponential delay
// between attempts. MaxAttempts includes the first request.
type Policy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	// Schedule optionally overrides exponential delays. It is useful when an
	// operational runbook calls for a non-power-of-two sequence such as
	// 30s/2m/10m.
	Schedule []time.Duration
}

// DefaultPolicy follows the operational plan: retry after roughly 30 seconds,
// then 2 minutes, then 10 minutes. The delays are capped so a malformed
// configuration cannot create an unbounded timer.
var DefaultPolicy = Policy{
	MaxAttempts:    4,
	InitialBackoff: 30 * time.Second,
	MaxBackoff:     10 * time.Minute,
	Schedule:       []time.Duration{30 * time.Second, 2 * time.Minute, 10 * time.Minute},
}

// Normalized returns a usable policy without mutating the caller's value.
func (policy Policy) Normalized() Policy {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.InitialBackoff < 0 {
		policy.InitialBackoff = 0
	}
	if policy.MaxBackoff <= 0 {
		policy.MaxBackoff = policy.InitialBackoff
	}
	if policy.MaxBackoff < policy.InitialBackoff {
		policy.MaxBackoff = policy.InitialBackoff
	}
	if policy.Schedule != nil {
		policy.Schedule = append([]time.Duration(nil), policy.Schedule...)
	}
	return policy
}

// Delay returns the delay before retryNumber (1 means the first retry).
func (policy Policy) Delay(retryNumber int) time.Duration {
	policy = policy.Normalized()
	if retryNumber <= 0 {
		return 0
	}
	if len(policy.Schedule) > 0 {
		index := retryNumber - 1
		if index >= len(policy.Schedule) {
			index = len(policy.Schedule) - 1
		}
		delay := policy.Schedule[index]
		if delay < 0 {
			return 0
		}
		if delay > policy.MaxBackoff {
			return policy.MaxBackoff
		}
		return delay
	}
	if policy.InitialBackoff == 0 {
		return 0
	}
	delay := policy.InitialBackoff
	for index := 1; index < retryNumber; index++ {
		if delay >= policy.MaxBackoff/2 {
			return policy.MaxBackoff
		}
		delay *= 2
	}
	if delay > policy.MaxBackoff {
		return policy.MaxBackoff
	}
	return delay
}

// Wait sleeps for delay, returning early when ctx is canceled.
func Wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

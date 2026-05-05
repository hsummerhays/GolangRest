package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
)

var (
	ErrExternalServiceUnavailable = errors.New("external pricing service is currently unavailable")
	ErrValidationFailed          = errors.New("pricing validation failed")
)

// PricingClient defines the interface for interacting with an external pricing service.
type PricingClient interface {
	ValidatePrice(ctx context.Context, name string, price float64) error
}

// MockPricingClient simulates an external HTTP-based pricing engine.
type MockPricingClient struct {
	MaxRetries int
	BaseDelay  time.Duration
}

func NewMockPricingClient() *MockPricingClient {
	return &MockPricingClient{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
	}
}

// ValidatePrice demonstrates professional resilience patterns: Context awareness,
// Retries with Exponential Backoff, and Jitter.
func (c *MockPricingClient) ValidatePrice(ctx context.Context, name string, price float64) error {
	var lastErr error

	for i := 0; i < c.MaxRetries; i++ {
		// Respect context cancellation/timeout during retries
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := c.doRequest(ctx, name, price)
		if err == nil {
			return nil // Success!
		}

		lastErr = err
		slog.Warn("External pricing call failed, retrying...", 
			"attempt", i+1, 
			"error", err,
			"product", name,
		)

		// Exponential Backoff with Jitter
		// This prevents "Thundering Herd" problems in distributed systems.
		delay := c.BaseDelay * (1 << i) // 2^i
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		
		time.Sleep(delay + jitter)
	}

	return fmt.Errorf("%w: %v", ErrExternalServiceUnavailable, lastErr)
}

// doRequest simulates a flaky network call to an external service.
func (c *MockPricingClient) doRequest(ctx context.Context, name string, price float64) error {
	// Simulate 30% failure rate to demonstrate our retry logic
	if rand.Float32() < 0.3 {
		return errors.New("network timeout")
	}

	// Simulate a 500 Internal Server Error from the external party
	if rand.Float32() < 0.1 {
		return errors.New("HTTP 500: Internal Server Error")
	}

	// Logic-based validation
	if price <= 0 {
		return ErrValidationFailed
	}

	return nil
}

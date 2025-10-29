package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100, // 100 requests
	window:   time.Minute,
}

// RateLimitInterceptor implements rate limiting per user/IP
func RateLimitInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Get user identifier (IP or user ID)
		identifier := getUserIdentifier(ctx)

		if !limiter.allow(identifier) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

func (rl *rateLimiter) allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get requests for this identifier
	requests := rl.requests[identifier]

	// Remove old requests
	valid := []time.Time{}
	for _, t := range requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	// Check if limit exceeded
	if len(valid) >= rl.limit {
		return false
	}

	// Add new request
	valid = append(valid, now)
	rl.requests[identifier] = valid

	return true
}

func getUserIdentifier(ctx context.Context) string {
	// Try to get user ID from context
	if userID, ok := ctx.Value("user_id").(int32); ok {
		return string(rune(userID))
	}

	// Otherwise use IP or default
	return "anonymous"
}

// Cleanup old entries periodically
func init() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-limiter.window)

			for key, requests := range limiter.requests {
				valid := []time.Time{}
				for _, t := range requests {
					if t.After(cutoff) {
						valid = append(valid, t)
					}
				}

				if len(valid) == 0 {
					delete(limiter.requests, key)
				} else {
					limiter.requests[key] = valid
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

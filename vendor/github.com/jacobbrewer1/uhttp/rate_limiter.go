package uhttp

import "golang.org/x/time/rate"

type RateLimiter interface {
	// Allow returns true if the request is allowed.
	Allow(key string) bool
}

type rateLimiter struct {
	// limiter is the limiter.
	limiters map[string]*rate.Limiter

	// rps is the requests per second.
	rps float64

	// burst is the burst.
	burst int
}

// NewRateLimiter creates a new rate limiter.
//
// rps is the requests per second.
func NewRateLimiter(rps float64) RateLimiter {
	return NewRateLimiterWithBurst(rps, 0)
}

// NewRateLimiterWithBurst creates a new rate limiter.
//
// rps is the requests per second.
//
// burst is the burst. This is the number of requests that can be made in one go. If the burst is 0, then the burst is
// set to the rps. If the burst is less than the rps, then the burst is set to the rps.
func NewRateLimiterWithBurst(rps float64, burst int) RateLimiter {
	if burst == 0 {
		burst = int(rps)
	} else if burst < int(rps) {
		burst = int(rps)
	}

	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

func (r *rateLimiter) Allow(key string) bool {
	// Rate limits the request.
	limiter, ok := r.limiters[key]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(r.rps), r.burst)
		r.limiters[key] = limiter
	}

	return limiter.Allow()
}

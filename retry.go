package xgo

import (
	"math"
	"math/rand"
	"time"
)

// Retrier decides whether to retry a failed function
type Retrier interface {
	Perform(func() error, func(err error) bool) error
}

// ExponentialBackoff represents retry with exponential backoff
type ExponentialBackoff struct {
	maxRetries      int
	maxRetrySeconds time.Duration
}

// NewExponentialBackoff creates a NewExponentialBackoff instance
func NewExponentialBackoff() *ExponentialBackoff {
	return &ExponentialBackoff{
		maxRetries:      10,
		maxRetrySeconds: 64 * time.Second,
	}
}

// Perform the exponential backoff algorithm
func (retrier *ExponentialBackoff) Perform(fn func() error, retryCondition func(err error) bool) error {
	population := int64(time.Second / time.Millisecond)

	var err error
	var backoff time.Duration
	var ms time.Duration
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < retrier.maxRetries; i++ {
		err = fn()
		if !retryCondition(err) {
			return nil
		}
		ms = time.Duration(r.Int63n(population)) * time.Millisecond
		backoff = time.Duration(math.Exp2(float64(i)))*time.Second + ms
		if backoff > retrier.maxRetrySeconds {
			backoff = retrier.maxRetrySeconds
		}
		time.Sleep(backoff)
	}

	return err
}

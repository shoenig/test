package wait

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeoutExceeded      = errors.New("wait: timeout exceeded")
	ErrAttemptsExceeded     = errors.New("wait: attempts exceeded")
	ErrConditionUnsatisfied = errors.New("wait: condition unsatisfied")
	ErrNoFunction           = errors.New("wait: no function specified")
)

const (
	defaultTimeout    = 3 * time.Second
	defaultGap        = 250 * time.Millisecond
	defaultIterations = 13
)

type Option func(*Context)

type runnable func(*runner) *result

type runner struct {
	ctx      *Context
	attempts int
}

type result struct {
	Err     error
	Success bool
}

// todo: Context parent

// Timeout sets the maximum amount of time to allow before giving up and marking
// the result as a failure.
//
// Default 3 seconds.
func Timeout(duration time.Duration) Option {
	return func(c *Context) {
		c.deadline = time.Now().Add(duration)
	}
}

// Attempts sets the maximum number of attempts to allow before giving up and
// marking the result as a failure.
//
// Default 12 attempts.
func Attempts(max int) Option {
	return func(c *Context) {
		c.iterations = max
	}
}

// Gap sets the amount of time to wait between attempts.
//
// Default 250 milliseconds.
func Gap(duration time.Duration) Option {
	return func(c *Context) {
		c.gap = duration
	}
}

// BoolFunc will retry f while it returns false, or a wait context threshold is
// exceeded.
func BoolFunc(f func() bool) Option {
	return func(c *Context) {
		c.r = boolFunc(f)
	}
}

func boolFunc(f func() bool) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctx.deadline)
		defer cancel()

		for {
			// make an attempt
			if f() {
				return &result{Success: true}
			}

			// used another attempt
			r.attempts++

			// check iterations
			if r.attempts > r.ctx.iterations {
				return &result{Err: ErrAttemptsExceeded}
			}

			// wait for gap or timeout
			select {
			case <-ctx.Done():
				return &result{Err: ErrTimeoutExceeded}
			case <-time.After(r.ctx.gap):
				// continue
			}
		}
	}
}

func ErrorFunc(f func() error) Option {
	return func(c *Context) {
		c.r = errorFunc(f)
	}
}

func errorFunc(f func() error) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctx.deadline)
		defer cancel()

		for {
			// make an attempt
			err := f()
			if err == nil {
				return &result{Success: true}
			}

			// used another attempt
			r.attempts++

			// check iterations
			if r.attempts > r.ctx.iterations {
				return &result{
					Err: fmt.Errorf("%v: %w", ErrAttemptsExceeded, err),
				}
			}

			// wait for gap or timeout
			select {
			case <-ctx.Done():
				return &result{
					Err: fmt.Errorf("%v: %w", ErrTimeoutExceeded, err),
				}
			case <-time.After(r.ctx.gap):
				// continue
			}
		}
	}
}

func TestFunc(f func() (bool, error)) Option {
	return func(c *Context) {
		c.r = testFunc(f)
	}
}

func testFunc(f func() (bool, error)) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctx.deadline)
		defer cancel()

		for {
			// make an attempt
			ok, err := f()
			if ok {
				return &result{Success: true}
			}

			// set default error
			if err == nil {
				err = ErrConditionUnsatisfied
			}

			// used another attempt
			r.attempts++

			// check iterations
			if r.attempts > r.ctx.iterations {
				return &result{
					Err: fmt.Errorf("%v: %w", ErrAttemptsExceeded, err),
				}
			}

			// wait for gap or timeout
			select {
			case <-ctx.Done():
				return &result{
					Err: fmt.Errorf("%v: %w", ErrTimeoutExceeded, err),
				}
			case <-time.After(r.ctx.gap):
				// continue
			}
		}
	}
}

func On(opts ...Option) *Context {
	c := &Context{now: time.Now()}
	for _, opt := range append([]Option{
		Timeout(defaultTimeout),
		Attempts(defaultIterations),
		Gap(defaultGap),
	}, opts...) {
		opt(c)
	}
	return c
}

type Context struct {
	now        time.Time
	deadline   time.Time
	gap        time.Duration
	iterations int
	r          runnable
}

func (ctx *Context) Run() error {
	if ctx.r == nil {
		return ErrNoFunction
	}
	return ctx.r(&runner{
		ctx:      ctx,
		attempts: 0,
	}).Err
}

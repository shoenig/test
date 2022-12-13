package wait

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrTimeoutExceeded      = errors.New("wait: timeout exceeded")
	ErrAttemptsExceeded     = errors.New("wait: attempts exceeded")
	ErrConditionUnsatisfied = errors.New("wait: condition unsatisfied")
	ErrNoFunction           = errors.New("wait: no function specified")
)

const (
	defaultTimeout = 3 * time.Second
	defaultGap     = 250 * time.Millisecond
)

type Option func(*Control)

type runnable func(*runner) *result

type runner struct {
	ctrl     *Control
	attempts int
}

type result struct {
	Err     error
	Success bool
}

// todo: Control parent

// Timeout sets the maximum amount of time to allow before giving up and marking
// the result as a failure.
//
// If set, the max attempts constraint is disabled.
//
// Default 3 seconds.
func Timeout(duration time.Duration) Option {
	return func(c *Control) {
		c.deadline = time.Now().Add(duration)
		c.iterations = math.MaxInt64
	}
}

// Attempts sets the maximum number of attempts to allow before giving up and
// marking the result as a failure.
//
// If set, the timeout constraint is disabled.
//
// By default a max timeout is used and the number of attempts is unlimited.
func Attempts(max int) Option {
	return func(c *Control) {
		c.iterations = max
		c.deadline = time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC)
	}
}

// Gap sets the amount of time to wait between attempts.
//
// Default 250 milliseconds.
func Gap(duration time.Duration) Option {
	return func(c *Control) {
		c.gap = duration
	}
}

// BoolFunc will retry f while it returns false, or a wait context threshold is
// exceeded.
func BoolFunc(f func() bool) Option {
	return func(c *Control) {
		c.r = boolFunc(f)
	}
}

func boolFunc(f func() bool) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctrl.deadline)
		defer cancel()

		for {
			// make an attempt
			if f() {
				return &result{Success: true}
			}

			// used another attempt
			r.attempts++

			// check iterations
			if r.attempts > r.ctrl.iterations {
				return &result{Err: ErrAttemptsExceeded}
			}

			// wait for gap or timeout
			select {
			case <-ctx.Done():
				return &result{Err: ErrTimeoutExceeded}
			case <-time.After(r.ctrl.gap):
				// continue
			}
		}
	}
}

func ErrorFunc(f func() error) Option {
	return func(c *Control) {
		c.r = errorFunc(f)
	}
}

func errorFunc(f func() error) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctrl.deadline)
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
			if r.attempts > r.ctrl.iterations {
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
			case <-time.After(r.ctrl.gap):
				// continue
			}
		}
	}
}

func TestFunc(f func() (bool, error)) Option {
	return func(c *Control) {
		c.r = testFunc(f)
	}
}

func testFunc(f func() (bool, error)) runnable {
	bg := context.Background()
	return func(r *runner) *result {
		ctx, cancel := context.WithDeadline(bg, r.ctrl.deadline)
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
			if r.attempts > r.ctrl.iterations {
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
			case <-time.After(r.ctrl.gap):
				// continue
			}
		}
	}
}

func On(opts ...Option) *Control {
	c := &Control{now: time.Now()}
	for _, opt := range append([]Option{
		Timeout(defaultTimeout),
		Gap(defaultGap),
	}, opts...) {
		opt(c)
	}
	return c
}

type Control struct {
	now        time.Time
	deadline   time.Time
	gap        time.Duration
	iterations int
	r          runnable
}

func (ctrl *Control) Run() error {
	if ctrl.r == nil {
		return ErrNoFunction
	}
	return ctrl.r(&runner{
		ctrl:     ctrl,
		attempts: 0,
	}).Err
}

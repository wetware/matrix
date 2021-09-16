package clock

import "time"

const defaultTimeStep = time.Millisecond

type Option func(c *Clock)

// WithTick sets the time-step for thes simulation clock.
// This is effectively precision with which the resulting
// clock is able to measure time.  A smaller time-step is
// more CPU-efficient.  If d < 0, defaults to millisecond
// precision.
//
// There is a trade-off between performance and precision
// such that a larger tick interval will reduce the load
// on the CPU when there are many events in the clock, at
// the expense of a
func WithTick(d time.Duration) Option {
	if d < 0 {
		d = defaultTimeStep
	}

	return func(c *Clock) {
		c.step = d
		c.ticks = ticker(d)
	}
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		WithTick(-1),
	}, opt...)
}

func ticker(accuracy time.Duration) func(time.Time) time.Duration {
	// native resolution?
	if accuracy == 0 || accuracy == 1 {
		return func(t time.Time) time.Duration {
			// avoid division operation for performance
			return time.Duration(t.UnixNano())
		}
	}

	return func(t time.Time) time.Duration {
		return time.Duration(t.UnixNano()) / accuracy
	}
}

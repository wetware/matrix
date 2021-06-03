package clock

import "time"

type Option func(c *Clock)

func WithAccuracy(d time.Duration) Option {
	if d < 0 {
		d = DefaultAccuracy
	}

	return func(c *Clock) {
		c.accuracy = d
		c.ticks = ticker(d)
	}
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		WithAccuracy(-1),
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

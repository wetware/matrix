package clock

import (
	"context"
	"math"
	"testing"
	"time"

	syncutil "github.com/lthibault/util/sync"
	"github.com/stretchr/testify/assert"
)

var t0 = time.Date(2021, 04, 9, 8, 0, 0, 0, time.UTC)

func TestTimeStep(t *testing.T) {
	t.Parallel()

	assert.Equal(t, defaultTimeStep,
		New().Timestep(), "wrong default time step")

	assert.Equal(t, time.Microsecond,
		New(WithTick(time.Microsecond)).Timestep(),
		"time step does not match option parameter")
}

func TestTicker(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("Accuracy=10ms", func(t *testing.T) {
		t.Parallel()

		c := New()

		var ctr syncutil.Ctr
		c.Ticker(time.Millisecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0) // init the clock
		c.Advance(t0.Add(time.Millisecond * 1010))

		assert.Eventually(t, func() bool {
			return ctr.Int() == 10
		}, time.Millisecond*100, time.Millisecond*10)
	})

	t.Run("Accuracy=10µs", func(t *testing.T) {
		t.Parallel()

		c := New(WithTick(time.Microsecond * 10))

		var ctr syncutil.Ctr
		c.Ticker(time.Microsecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0) // init the clock
		c.Advance(t0.Add(time.Microsecond * 1010))

		assert.Eventually(t, func() bool {
			return ctr.Int() == 10
		}, time.Millisecond*100, time.Millisecond*10)
	})
}

func TestMaxVal(t *testing.T) {
	t.Parallel()

	assert.Equal(t, maxVal(), uint64(math.MaxUint32))
}

func TestLevelMax(t *testing.T) {
	t.Parallel()

	assert.Equal(t, levelMax(1), uint64(1<<(nearShift+levelShift)))
	assert.Equal(t, levelMax(2), uint64(1<<(nearShift+2*levelShift)))
	assert.Equal(t, levelMax(3), uint64(1<<(nearShift+3*levelShift)))
	assert.Equal(t, levelMax(4), uint64(1<<(nearShift+4*levelShift)))
}

func TestGenVersion(t *testing.T) {
	t.Parallel()

	assert.Equal(t, genVersionHeight(1, 0xf), uint64(0x0001000f00000000))
	assert.Equal(t, genVersionHeight(1, 64), uint64(0x0001004000000000))
}

func TestHour(t *testing.T) {
	t.Parallel()

	c := New(incr)

	testHour := new(bool)
	done := make(chan struct{}, 1)
	c.After(time.Hour, func() {
		*testHour = true
		done <- struct{}{}
	})

	expire := c.getExpire(time.Hour, 0)
	for i := 0; i < int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-done:
	case <-time.After(time.Second / 100):
	}
	assert.True(t, *testHour)
}

func TestTicker_5s(t *testing.T) {
	t.Parallel()

	const total = 1000

	var (
		c        = New(incr)
		first5   syncutil.Ctr
		testTime = time.Second * 5
	)

	ctx, cancel := context.WithCancel(context.Background())
	c.Ticker(testTime, func() {
		if first5.Incr() == total {
			cancel()
		}
	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5.Int())
}

func TestTicker_hour(t *testing.T) {
	t.Parallel()

	const total = 100

	var (
		c        = New(incr)
		first5   syncutil.Ctr
		testTime = time.Hour
	)

	ctx, cancel := context.WithCancel(context.Background())
	c.Ticker(testTime, func() {
		if first5.Incr() == total {
			cancel()
		}
	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5.Int())
}

func TestTicker_day(t *testing.T) {
	t.Parallel()

	const total = 10

	var (
		c        = New(incr)
		first5   syncutil.Ctr
		testTime = time.Hour * 24
	)

	ctx, cancel := context.WithCancel(context.Background())
	c.Ticker(testTime, func() {
		if first5.Incr() == total {
			cancel()
		}
	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5.Int())
}

func TestAfter(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("Accuracy=10ms", func(t *testing.T) {
		t.Parallel()

		c := New()

		var ctr syncutil.Ctr
		c.After(time.Millisecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0)
		c.Advance(t0.Add(time.Second))

		assert.Eventually(t, func() bool {
			return ctr.Int() == 1
		}, time.Millisecond*100, time.Millisecond*10)
	})

	t.Run("Accuracy=µs", func(t *testing.T) {
		t.Parallel()

		c := New(WithTick(time.Microsecond * 10))

		var ctr syncutil.Ctr
		c.After(time.Microsecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0)
		c.Advance(t0.Add(time.Second))

		assert.Eventually(t, func() bool {
			return ctr.Int() == 1
		}, time.Millisecond*100, time.Millisecond*10)
	})
}

func TestNodeStop_1(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const d = time.Millisecond * 10

	var (
		c     = New()
		count syncutil.Ctr
	)

	stop := c.After(d, func() { count.Incr() })

	go func() {
		time.Sleep(d * 2)
		stop()
		cancel()
	}()

	run(ctx, c, d/2)
	assert.NotEqual(t, count.Int(), 1)
}

func TestNodeStop(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const d = time.Millisecond * 100

	var (
		c     = New()
		count syncutil.Ctr
	)

	stop := c.After(d, func() { count.Incr() })
	stop()
	go func() {
		time.Sleep(d * 2)
		cancel()
	}()

	run(ctx, c, d/2)
	assert.NotEqual(t, count, 1)
}

func run(ctx context.Context, c *Clock, d time.Duration) {
	var t = t0
	for ctx.Err() == nil {
		c.Advance(t)
		t = t0.Add(d)
	}

	// tk := time.NewTicker(d)
	// defer tk.Stop()

	// for {
	// 	select {
	// 	case t := <-tk.C:
	// 		c.Advance(t)
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

func incr(c *Clock) {
	WithTick(-1)(c)
	c.ticks = func(time.Time) time.Duration {
		return c.curTimePoint + 1
	}
}

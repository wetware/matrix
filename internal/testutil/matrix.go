package testutil

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
	mock_mx "github.com/wetware/matrix/internal/mock/pkg"
)

func NewHostFactory(ctrl *gomock.Controller) *mock_mx.MockHostFactory {
	var ctx = reflect.TypeOf((*context.Context)(nil)).Elem()

	f := mock_mx.NewMockHostFactory(ctrl)
	f.EXPECT().
		NewHost(
			gomock.AssignableToTypeOf(ctx),
			gomock.AssignableToTypeOf([]config.Option(nil)),
		).
		DoAndReturn(func(ctx context.Context, opt []config.Option) (host.Host, error) {
			return NewHost(ctrl), nil
		}).
		AnyTimes()

	return f
}

func NewClock(ctrl *gomock.Controller, accuracy time.Duration, onTick func(t time.Time)) *mock_mx.MockClockController {
	if onTick == nil {
		onTick = func(time.Time) {}
	}

	c := mock_mx.NewMockClockController(ctrl)
	c.EXPECT().
		Advance(&monotonicTimeMatcher{}).
		Do(onTick).
		AnyTimes()

	if accuracy == 0 {
		accuracy = time.Millisecond * 10
	}

	c.EXPECT().
		Accuracy().
		Return(accuracy).
		AnyTimes()

	return c
}

type monotonicTimeMatcher struct{ time.Time }

// Matches returns whether x is a match.
func (m *monotonicTimeMatcher) Matches(x interface{}) bool {
	if t, ok := x.(time.Time); ok && t.After(m.Time) {
		m.Time = t
		return true
	}

	return false
}

// String describes what the matcher matches.
func (m *monotonicTimeMatcher) String() string { return fmt.Sprintf("time after %s", m.Time) }

package testutil

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	mock_mx "github.com/wetware/matrix/internal/mock/pkg/matrix"
)

func NewClock(ctrl *gomock.Controller, accuracy time.Duration, onTick func(t time.Time)) *mock_mx.MockClockController {
	c := mock_mx.NewMockClockController(ctrl)

	if onTick != nil {
		c.EXPECT().
			Advance(&monotonicTimeMatcher{}).
			Do(onTick).
			AnyTimes()
	}

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

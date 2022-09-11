package date_test

import (
	"testing"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/date"
	"github.com/stretchr/testify/assert"
)

func TestTruncateToDay(t *testing.T) {
	now := time.Now()
	truncated := date.TruncateToDay(now)
	assert.Equal(t, now.Year(), truncated.Year())
	assert.Equal(t, now.Month(), truncated.Month())
	assert.Equal(t, now.Day(), truncated.Day())
	assert.Equal(t, 0, truncated.Hour(), "hour should truncated to 0")
	assert.Equal(t, 0, truncated.Minute(), "minute should truncated to 0")
	assert.Equal(t, 0, truncated.Second(), "second should truncated to 0")
	assert.Equal(t, 0, truncated.Nanosecond(), "nano second should truncated to 0")
}

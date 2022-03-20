package util

import (
	"fmt"
	"time"
)

func Timer() func() string {
	start := time.Now()
	return func() string {
		return formatDurationSec(time.Now().Sub(start))
	}
}

func formatDurationSec(duration time.Duration) string {
	return fmt.Sprintf("%.2f", duration.Seconds())
}

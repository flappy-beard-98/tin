package core

import (
	"testing"
)

func TestComplexPercentRising(t *testing.T) {
	value, multiplier, rounds := 1.0, 1.0, 365

	for i := 0; i < rounds; i++ {
		value += (float64(value) / 100.0) * multiplier
	}

	t.Logf("for %.0fx multiplier and %d rounds value = %.2f", multiplier, rounds, value)
}

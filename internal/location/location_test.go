package location

import (
	"testing"
)

func TestForIP(t *testing.T) {
	const ip string = "24.48.0.1"

	t.Run("runs successfully", func(t *testing.T) {
		if _, err := ForIP(ip); err != nil {
			t.Errorf("%v", err.Error())
		}
	})
}

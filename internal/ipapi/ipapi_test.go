package ipapi

import (
	"testing"
)

func TestLocateIP(t *testing.T) {
	const ip string = "24.48.0.1"

	t.Run("runs successfully", func(t *testing.T) {
		if _, err := LocateIP(ip); err != nil {
			t.Errorf("%v", err.Error())
		}
	})
}

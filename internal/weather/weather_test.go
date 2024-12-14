package weather

import (
	"testing"
)

func Test(t *testing.T) {
	const lat float64 = 8.716667
	const lon float64 = 167.733333

	t.Run("runs successfully", func(t *testing.T) {
		if _, err := ForLatLon(lat, lon); err != nil {
			t.Errorf("%v", err.Error())
		}
	})
}

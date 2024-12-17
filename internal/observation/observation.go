package observation

import (
	"weather/internal/data"

	"time"
)

func ResolvePriorObservation(obs data.Observation, db *data.Queries) (*data.Observation, error) {
	return &data.Observation{
		ID:               0,
		Latitude:         0.00,
		Longitude:        0.00,
		Timezone:         "GMT",
		TempC:            0.0,
		TempF:            32.0,
		RelativeHumidity: 0.0,
		Rain:             0.0,
		Snowfall:         0.0,
		WeatherCode:      "2",
		TimeUtc:          time.Now().UTC(),
		TimeLocal:        time.Now().UTC(),
	}, nil
}

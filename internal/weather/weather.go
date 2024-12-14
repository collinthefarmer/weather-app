package weather

import (
	"weather/internal/fetch"
	"weather/internal/validation"

	"fmt"
)

type CurrentUnits struct {
	Time               string `json:"time"`
	Interval           string `json:"interval"`
	Temperature2m      string `json:"temperature_2m"`
	RelativeHumidity2m string `json:"relative_humidity_2m"`
	Rain               string `json:"rain"`
	Showers            string `json:"showers"`
	Snowfall           string `json:"snowfall"`
}

type Current struct {
	Time               string  `json:"time"`
	Interval           int     `json:"interval"`
	Temperature2m      float64 `json:"temperature_2m"`
	RelativeHumidity2m int     `json:"relative_humidity_2m"`
	Rain               float64 `json:"rain"`
	Showers            float64 `json:"showers"`
	Snowfall           float64 `json:"snowfall"`
}

type Weather struct {
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	GenerationTimeMs     float64      `json:"generation_time_ms"`
	UTCOffsetSeconds     int          `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	CurrentUnits         CurrentUnits `json:"current_units"`
	Current              Current      `json:"current"`
}

func (w Weather) Validate() (validation.ValidationProblems, error) {
	// TODO
	return nil, nil
}

const basePath = "https://api.open-meteo.com/v1/forecast"
const fields = "temperature_2m,relative_humidity_2m,rain,showers,snowfall"

func ForLatLon(lat float64, lon float64) (Weather, error) {
	weather := Weather{}

	endpoint := fmt.Sprintf("%s?current=%s&latitude=%.2f&longitude=%.2f",
		basePath, fields, lat, lon,
	)

	if err := fetch.JSON(endpoint, &weather); err != nil {
		return weather, fmt.Errorf("OpenMeteo API error %w", err)
	}

	return weather, nil
}

package display

import (
	"weather/internal/location"
	"weather/internal/weather"

	"html/template"
	"net/http"
)

type Data struct {
	Location location.Geolocation
	Weather  weather.Weather
}

func SelectTemplate(loc location.Geolocation, wth weather.Weather) (*template.Template, error) {
	return nil, nil
}

func RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data *Data) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}

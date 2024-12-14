package main

import (
	"weather/internal/data"
	"weather/internal/location"
	"weather/internal/weather"

	"embed"
	"html/template"
	"net/http"
	"strings"
)

func renderTemplate[T any](w http.ResponseWriter, tmpl *template.Template, data *T) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
	return tmpl.Execute(w, data)
}

func HandleIndexGet(tmpl *template.Template, db *data.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestAddress := r.RemoteAddr
		requestIP := strings.Split(requestAddress, ":")[0]

		// figure out where the request is coming from
		loc, err := location.ForIP(requestIP)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your location :("))
		}

		entry, err := db.AddGeolocation(ctx, data.AddGeolocationParams{
			Ip:        requestIP,
			Latitude:  loc.Lat,
			Longitude: loc.Lon,
		})

		// find the most recent weather for the location of the request
		wth, err := weather.ForLatLon(loc.Lat, loc.Lon)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find you weather :("))
		}

		bestImage := "" // as determined by some math? closest + most recent - not sure how to do this

		data := struct {
			Location     location.Geolocation
			Weather      weather.Weather
			ImageDataUri string
		}{Location: loc, Weather: wth}

		// present data to user
		// alongside image of observation
		// prompt to draw what they see

		// JS canvas
		if err := renderTemplate(w, tmpl, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I beefed it"))
		}
	})
}

func HandleIndexPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {

	})
}

//go:embed templates/*
var templateFS embed.FS

func main() {
	templates, err := template.ParseFS(templateFS, "templates/*.template.html")
	if err != nil {
		panic(err)
	}

	server := http.NewServeMux()

	server.Handle("GET /", HandleIndexGet(templates))
	server.Handle("POST /", HandleIndexPost())

	http.ListenAndServe("localhost:8080", server)
}

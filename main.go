package main

import (
	"weather/internal/display"
	"weather/internal/location"
	"weather/internal/weather"

	"net/http"
	"strings"
)

func HandleIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAddress := r.RemoteAddr
		requestIP := strings.Split(requestAddress, ":")[0]

		// figure out where the request is coming from
		loc, err := location.ForIP(requestIP)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your location :("))
		}

		// find the most recent weather for the location of the request
		wth, err := weather.ForLatLon(loc.Lat, loc.Lon)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find you weather :("))
		}

		data := display.Data{Location: loc, Weather: wth}

		// return a template corresponding to the weather of that location
		tmpl, err := display.SelectTemplate(loc, wth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I beefed it"))
		}

		if err := display.RenderTemplate(w, tmpl, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I beefed it"))
		}
	})
}

func main() {
	server := http.NewServeMux()

	server.Handle("/", HandleIndex())

	http.ListenAndServe("localhost:8080", server)
}

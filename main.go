package main

import (
	"log"
	"weather/internal/data"
	"weather/internal/drawing"
	"weather/internal/location"
	"weather/internal/validation"
	"weather/internal/weather"

	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getGeolocation(r *http.Request, db *data.Queries) (data.Geolocation, error) {
	ctx := r.Context()

	requestAddress := r.RemoteAddr
	ip := strings.Split(requestAddress, ":")[0]

	entry, err := db.GetGeolocation(ctx, ip)
	switch err {
	case nil:
		return entry, nil
	case sql.ErrNoRows:
		loc, err := location.ForIP(ip)
		if err != nil {
			return entry, err
		}

		entry, err = db.AddGeolocation(ctx, data.AddGeolocationParams{
			Ip:        ip,
			Latitude:  loc.Lat,
			Longitude: loc.Lon,
		})

		return entry, nil
	default:
		return entry, err
	}
}

func renderTemplate[T any](w http.ResponseWriter, tmpl *template.Template, data *T, status int) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/html")
	return tmpl.Execute(w, data)
}

func handleIndexGet(tmpl *template.Template, db *data.Queries) http.Handler {
	const indexTemplateName = "index"
	indexTemplate := tmpl.Lookup("index")
	if indexTemplate == nil {
		panic("")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc, err := getGeolocation(r, db)
		if err != nil {
			log.Printf("%v", err)

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your location :("))
			return
		}

		wth, err := weather.ForLatLon(loc.Latitude, loc.Longitude)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your weather :("))
			return
		}

		bestImage := "" // as determined by some math? closest + most recent - not sure how to do this

		data := struct {
			Location     data.Geolocation
			Weather      weather.Weather
			ImageDataUri string
		}{Location: loc, Weather: wth, ImageDataUri: bestImage}

		if err := renderTemplate(w, tmpl, &data, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I beefed it"))
			return
		}
	})
}

func readObservationDrawing(r *http.Request) (*data.AddObservationDrawingParams, error) {
	id := r.PathValue("id")
	if id == "" {
		return nil, validation.ErrValidation
	}

	drawingData := r.PostFormValue("drawing")
	err := drawing.Validate(drawingData)
	if err != nil {
		return nil, validation.ErrValidation
	}

	return &data.AddObservationDrawingParams{
		DrawingDataUri:   sql.NullString{String: drawingData, Valid: true},
		DrawingSizeBytes: sql.NullInt64{Int64: int64(len([]byte(drawingData)))},
	}, nil
}

func handleObservationPatch(templates *template.Template, db *data.Queries) http.Handler {
	const observationTemplateName = "observation"
	observationTemplate := templates.Lookup(observationTemplateName)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		obs, err := readObservationDrawing(r)
		switch err {
		case nil:
			err := db.AddObservationDrawing(ctx, *obs)
			if err == sql.ErrNoRows {
				http.Error(w, "", http.StatusNotFound)
				return
			} else if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			renderTemplate(w, observationTemplate, obs, http.StatusCreated)
			return
		case validation.ErrValidation:
			http.Error(w, "", http.StatusBadRequest)
			return
		default:
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}

func createDatabase(path string, ddl string) (*data.Queries, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("couldn't open database connection: %w", err)
	}

	_, err = db.Exec(ddl)
	if err != nil {
		return nil, fmt.Errorf("could apply database schema: %w", err)
	}

	return data.New(db), nil
}

//go:embed sqlite/schema.sql
var schema string

//go:embed templates/*
var templateFS embed.FS

func main() {
	templates, err := template.ParseFS(templateFS, "templates/*.template.html")
	if err != nil {
		panic(err)
	}

	db, err := createDatabase("./db.sqlite", schema)
	if err != nil {
		panic(err)
	}

	server := http.NewServeMux()

	server.Handle("GET /", handleIndexGet(templates, db))
	server.Handle("PATCH /observations/{id}", handleObservationPatch(templates, db))

	http.ListenAndServe("localhost:8080", server)
}

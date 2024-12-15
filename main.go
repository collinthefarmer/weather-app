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
		log.Printf("fetching location for %v", ip)

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

func renderTemplate[T any](w http.ResponseWriter, tmpl *template.Template, data T) error {
	err := tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

type indexTemplateData struct {
	IP        string
	Latitude  float64
	Longitude float64
}

func handleIndexGet(tmpl *template.Template, db *data.Queries) http.Handler {
	const indexTemplateName = "index"
	indexTemplate := tmpl.Lookup("index")
	if indexTemplate == nil {
		panic("no index template")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc, err := getGeolocation(r, db)
		if err != nil {
			log.Printf("error getting geolocation: %v", err)

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your location :("))
			return
		}

		wth, err := weather.ForLatLon(loc.Latitude, loc.Longitude)
		if err != nil {
			log.Printf("error getting weather: %v", err)

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I couldn't find your weather :("))
			return
		}

		bestImage := "" // as determined by some math? closest + most recent - not sure how to do this
		_, _ = bestImage, wth

		data := indexTemplateData{
			IP:        loc.Ip,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		}

		if err := renderTemplate(w, tmpl, data); err != nil {
			panic(err)
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
	if observationTemplate == nil {
		panic("no observation template")
	}

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

			renderTemplate(w, observationTemplate, *obs)
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
	templates, err := template.ParseFS(
		templateFS,
		"templates/*.template.html",
		"templates/js/*.template.html",
	)
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

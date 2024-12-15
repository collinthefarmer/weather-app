package main

import (
	"weather/internal/data"
	"weather/internal/drawing"
	"weather/internal/location"
	"weather/internal/observation"
	"weather/internal/validation"
	"weather/internal/weather"

	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func resolveGeolocation(ctx context.Context, ip string, db *data.Queries) (data.Geolocation, error) {

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
			City:      loc.City,
			Country:   loc.Country,
			Timezone:  loc.Timezone,
		})

		return entry, nil
	default:
		return entry, err
	}
}

func resolveObservation(ctx context.Context, loc data.Geolocation, db *data.Queries) (*data.Observation, error) {
	// todo: make this work like resolveGeolocation
	// then refactor, probably
	wth, err := weather.ForLatLon(loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, err
	}

	obs, err := db.AddObservation(ctx, data.AddObservationParams{
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timezone:  loc.Timezone,
		TempC:     wth.Current.Temperature2m,
	})
	if err != nil {
		return nil, err
	}

	return &obs, nil
}

func renderTemplate[T any](w http.ResponseWriter, tmpl *template.Template, data T) error {
	err := tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

type indexTemplateData struct {
	PrevObservation data.Observation
	NextObservation data.Observation
}

func handleIndexGet(tmpl *template.Template, db *data.Queries) http.Handler {
	const indexTemplateName = "index"
	indexTemplate := tmpl.Lookup("index")
	if indexTemplate == nil {
		panic("no index template")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestAddress := r.RemoteAddr
		ip := strings.Split(requestAddress, ":")[0]

		loc, err := resolveGeolocation(ctx, ip, db)
		if err != nil {
			log.Printf("error resolving geolocation: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I couldn't find your location :("))
			return
		}

		obs, err := resolveObservation(ctx, loc, db)
		if err != nil {
			log.Printf("error resolving observation: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("uh oh, I couldn't find your weather :("))
			return
		}

		prev, err := observation.ResolvePriorObservation(*obs, db)
		if err != nil {
			log.Printf("error resolving previous observation: %v", err)

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("uh oh, I beefed it!"))
			return
		}

		if err := renderTemplate(w, tmpl, indexTemplateData{
			PrevObservation: *prev,
			NextObservation: *obs,
		}); err != nil {
			log.Printf("error rendering index template: %v", err)
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
		log.Fatalf("error parsing templates: %v", err)
	}

	db, err := createDatabase("./db.sqlite", schema)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}

	server := http.NewServeMux()

	server.Handle("GET /", handleIndexGet(templates, db))
	server.Handle("PATCH /observations/{id}", handleObservationPatch(templates, db))

	http.ListenAndServe("localhost:8080", server)
}

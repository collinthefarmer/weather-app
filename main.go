package main

import (
	"weather/internal/data"
	"weather/internal/drawing"
	"weather/internal/location"
	"weather/internal/observation"
	"weather/internal/templates"
	"weather/internal/validation"
	"weather/internal/weather"

	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		if err != nil {
			return entry, nil
		}

		return entry, nil
	default:
		return entry, err
	}
}

func resolveObservation(ctx context.Context, loc data.Geolocation, db *data.Queries) (*data.Observation, error) {
	wth, err := weather.ForLatLon(loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, err
	}

	tzloc, err := time.LoadLocation(loc.Timezone)
	if err != nil {
		tzloc = time.UTC
	}

	obs, err := db.AddObservation(ctx, data.AddObservationParams{
		Latitude:         loc.Latitude,
		Longitude:        loc.Longitude,
		Timezone:         loc.Timezone,
		TempC:            wth.Current.Temperature2m,
		TempF:            wth.Current.Temperature2m,
		Rain:             wth.Current.Rain,
		Snowfall:         wth.Current.Snowfall,
		WeatherCode:      strconv.Itoa(wth.Current.WeatherCode),
		RelativeHumidity: float64(wth.Current.RelativeHumidity2m),
		TimeUtc:          time.Now().UTC(),
		TimeLocal:        time.Now().In(tzloc),
	})

	if err != nil {
		return nil, err
	}

	return &obs, nil
}

func resolveObservationByID(ctx context.Context, id int64) (*data.Observation, error) {
	// TODO
	return nil, nil
}

func readObservationDrawing(r *http.Request) (*data.ObservationDrawing, error) {
	drawingData := r.PostFormValue("drawing")
	idStr := r.PathValue("id")
	if idStr == "" {
		return nil, validation.ErrValidation
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, validation.ErrValidation
	}

	if err := drawing.Validate(drawingData); err != nil {
		return nil, validation.ErrValidation
	}

	return &data.ObservationDrawing{
		ObservationID: int64(id),
		Data:          drawingData,
		SizeBytes:     int64(len([]byte(drawingData))),
		TimeSubmitted: time.Now().UTC(),
	}, nil
}

func createObservationDrawing(ctx context.Context, drawing *data.ObservationDrawing, db *data.Queries) error {
	return db.AddObservationDrawing(ctx, data.AddObservationDrawingParams{
		ObservationID: drawing.ObservationID,
		Data:          drawing.Data,
		SizeBytes:     drawing.SizeBytes,
		TimeSubmitted: drawing.TimeSubmitted,
	})
}

func handleIndexGet(tmpl *templates.TemplateEngine, db *data.Queries) http.Handler {
	const indexTemplateName = "templates/index.template.html"

	type indexTemplateData struct {
		Location        data.Geolocation
		PrevObservation data.Observation
		NextObservation data.Observation
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestAddress := r.RemoteAddr
		ip := strings.Split(requestAddress, ":")[0]

		loc, err := resolveGeolocation(ctx, ip, db)
		if err != nil {
			switch err {
			case context.Canceled:
			default:
				log.Printf("error resolving geolocation: %v", err)
				break
			}

			http.Error(w, "uh oh, I couldn't find your location :(", http.StatusInternalServerError)
			return
		}

		obs, err := resolveObservation(ctx, loc, db)
		if err != nil {
			switch err {
			case context.Canceled:
			default:
				log.Printf("error resolving observation: %v", err)
				break
			}

			http.Error(w, "uh oh, I couldn't find your weather :(", http.StatusInternalServerError)
			return
		}

		prev, err := observation.ResolvePriorObservation(ctx, *obs, db)
		if err != nil {
			switch err {
			case context.Canceled:
			default:
				log.Printf("error resolving prior observation: %v", err)
				break
			}

			http.Error(w, "uh oh, I beefed it :(", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Render(w, indexTemplateName, indexTemplateData{
			Location:        loc,
			PrevObservation: *prev,
			NextObservation: *obs,
		}); err != nil {
			log.Printf("error rendering index template: %v", err)
			return
		}
	})
}

func handleObservationDrawingPost(tmpl *templates.TemplateEngine, db *data.Queries) http.Handler {
	const observationTemplateName = "templates/fragments/observation.template.html"

	type observationTemplateData struct {
		Observation data.Observation
		Drawing     data.ObservationDrawing
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		drawing, err := readObservationDrawing(r)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		observation, err := db.GetObservation(ctx, drawing.ObservationID)
		switch err {
		case nil:
			break
		case sql.ErrNoRows: // TODO: figure out what kind of error is returned on foreign key constraint failed
			http.Error(w, "", http.StatusBadRequest)
			return
		default:
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := createObservationDrawing(ctx, drawing, db); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

		tmpl.Render(w, observationTemplateName, observationTemplateData{
			Observation: observation,
			Drawing:     *drawing,
		})
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

//go:embed static/*
var staticFS embed.FS

var templateConstants = struct {
	MinLatitude  float32
	MaxLatitude  float32
	MinLongitude float32
	MaxLongitude float32
}{
	MinLatitude:  -90.0,
	MaxLatitude:  90.0,
	MinLongitude: -180.0,
	MaxLongitude: 180.0,
}

func main() {
	templates, err := templates.Init(
		templateFS,
		templateConstants,
		"templates/root.template.html",
		"templates/common/*.template.html",
		"templates/fragments/*.template.html",
	)
	if err != nil {
		log.Fatalf("error parsing templates: %v", err)
	}

	db, err := createDatabase("./db.sqlite", schema)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}

	server := http.NewServeMux()

	server.Handle(
		"GET /static/",
		http.FileServerFS(staticFS),
	)

	server.Handle(
		"GET /",
		handleIndexGet(templates, db),
	)

	server.Handle(
		"POST /observations/{id}/drawings",
		handleObservationDrawingPost(templates, db),
	)

	http.ListenAndServe("localhost:8080", server)
}

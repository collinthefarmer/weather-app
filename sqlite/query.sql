-- name: AddGeolocation :one
INSERT INTO
    geolocations (ip, latitude, longitude, city, country, timezone)
VALUES
    (?, ?, ?, ?, ?, ?)
RETURNING
    *;

-- name: GetGeolocation :one
SELECT
    *
FROM
    geolocations
WHERE
    ip = ?;

-- name: AddObservation :one
INSERT INTO
    observations (
        latitude,
        longitude,
        timezone,
        temp_c,
        temp_f,
        relative_humidity,
        rain,
        snowfall,
        weather_code,
        time_utc,
        time_local
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING
    *;

-- name: AddObservationDrawing :exec
INSERT INTO
    observations_drawings (data, size_bytes, time_submitted)
SET
    drawing_data_uri = ?,
    drawing_size_bytes = ?,
    time_drawing_submit = ?
RETURNING
    *;

-- todo
-- name: SelectRecentObservation :one
SELECT
    *
FROM
    observations
LIMIT
    1;

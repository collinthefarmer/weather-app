-- name: AddGeolocation :one
INSERT INTO
    geolocations (ip, latitude, longitude)
VALUES
    (?, ?, ?)
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
        time_utc,
        time_local,
        timezone,
        latitude,
        longitude,
        temperature_2m,
        relative_humidity_2m,
        rain,
        showers,
        snowfall,
        drawing_data_uri,
        drawing_size_bytes
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING
    *;

-- name: AddObservationDrawing :exec
UPDATE
    observations
SET
    drawing_data_uri = ?,
    drawing_size_bytes = ?
WHERE
    id = ?
    AND drawing_data_uri IS NULL;

-- todo
-- name: SelectRecentObservation :one
SELECT
    *
FROM
    observations
LIMIT
    1;

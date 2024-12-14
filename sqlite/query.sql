-- name: AddGeolocation :one
INSERT INTO
    geolocations (ip, latitude, longitude)
VALUES
    (?, ?, ?)
RETURNING
    *;

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

-- todo
-- name: SelectRecentObservation :one
SELECT
    *
FROM
    observations
LIMIT
    1;

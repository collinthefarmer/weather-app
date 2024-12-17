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

-- name: GetObservation :one
SELECT
    *
FROM
    observations
WHERE
    id = ?;

-- name: AddObservationDrawing :exec
INSERT INTO
    observation_drawings (observation_id, data, size_bytes, time_submitted)
VALUES
    (?, ?, ?, ?)
RETURNING
    *;

-- todo - returns best
-- name: PriorObservation :one
SELECT
    *
FROM
    observations o
    INNER JOIN observation_drawings od ON o.id = od.observation_id;

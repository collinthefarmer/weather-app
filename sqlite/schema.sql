CREATE TABLE IF NOT EXISTS geolocations (
    ip TEXT PRIMARY KEY,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    city TEXT NOT NULL,
    country TEXT NOT NULL,
    timezone TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS observations (
    id INTEGER PRIMARY KEY,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    timezone TEXT NOT NULL,
    temp_c REAL NOT NULL,
    temp_f REAL NOT NULL,
    relative_humidity REAL NOT NULL,
    rain REAL NOT NULL,
    snowfall REAL NOT NULL,
    weather_code TEXT NOT NULL,
    time_utc DATETIME NOT NULL,
    time_local DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS observation_drawings (
    data TEXT NOT NULL,
    size_bytes INT NOT NULL,
    time_submitted DATETIME NOT NULL
);

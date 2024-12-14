CREATE TABLE IF NOT EXISTS geolocations (
    ip TEXT PRIMARY KEY,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS observations (
    id INTEGER PRIMARY KEY,
    time_utc DATETIME NOT NULL,
    time_local DATETIME NOT NULL,
    timezone TEXT NOT NULL,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    temperature_2m REAL NOT NULL,
    relative_humidity_2m REAL NOT NULL,
    rain REAL NOT NULL,
    showers REAL NOT NULL,
    snowfall REAL NOT NULL,
    drawing_data_uri TEXT,
    drawing_size_bytes INTEGER
);

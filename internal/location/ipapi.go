package location

import (
	"fmt"
	"net/url"
	"weather/internal/fetch"
	"weather/internal/validation"
)

type IPAPIGeolocation struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
}

func (loc IPAPIGeolocation) Validate() (validation.ValidationProblems, error) {
	// TODO
	return nil, nil
}

const basePath = "http://ip-api.com/json/"
const fields = "status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,query"
const defaultIP = "127.0.0.1"

func ForIP(ip string) (IPAPIGeolocation, error) {
	geolocation := IPAPIGeolocation{}

	endpoint := basePath
	if ip == defaultIP {
		ip = ""
	}

	endpoint, err := url.JoinPath(basePath, ip)
	if err != nil {
		return geolocation, fmt.Errorf("error building IP-API path for ip: %s: %w", ip, err)
	}

	if err := fetch.JSON(endpoint, &geolocation); err != nil {
		return geolocation, fmt.Errorf("error communicating with IP-API.com, %w", err)
	}

	return geolocation, nil
}

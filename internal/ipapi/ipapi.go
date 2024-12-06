package ipapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Geolocation struct {
	Query       string
	Status      string
	Country     string
	CountryCode string
	Region      string
	RegionName  string
	City        string
	Zip         string
	Lat         float32
	Lon         float32
	Timezone    string
}

func (loc Geolocation) Validate() error {
	// surely there has GOT to be a better way of doing with w/o Reflect...

	if loc.Timezone == "" {
		return errors.New("missing value for Timezone")
	}
	if loc.Lon == 0. {
		return errors.New("missing value for Lon")
	}
	if loc.Lat == 0. {
		return errors.New("missing value for Lat")
	}
	if loc.Zip == "" {
		return errors.New("missing value for Zip")
	}
	if loc.City == "" {
		return errors.New("missing value for City")
	}
	if loc.RegionName == "" {
		return errors.New("missing value for RegionName")
	}
	if loc.Region == "" {
		return errors.New("missing value for Region")
	}
	if loc.CountryCode == "" {
		return errors.New("missing value for CountryCode")
	}
	if loc.Country == "" {
		return errors.New("missing value for Country")
	}
	if loc.Status == "" {
		return errors.New("missing value for Status")
	}
	if loc.Query == "" {
		return errors.New("missing value for Query")
	}
	return nil
}

const basePath = "http://ip-api.com/json/"
const fields = "status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,query"

func fetchFromIPAPI(ip string, geolocation *Geolocation) error {
	if ip == "127.0.0.1" {
		ip = ""
	}

	apiPath, err := url.JoinPath(basePath, ip)
	if err != nil {
		return fmt.Errorf("problem constructing IP-API path for IP %s: %w", ip, err)
	}

	apiPath += "?fields=" + fields

	response, err := http.Get(apiPath)
	if err != nil {
		return fmt.Errorf("error contacting IP-API: %w", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("non-200 status code returned from IP-API: %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading IP-API response body: %w", err)
	}

	if err := json.Unmarshal(body, &geolocation); err != nil {
		return fmt.Errorf("error decodings JSON from IP-API response body: %w", err)
	}

	if err := geolocation.Validate(); err != nil {
		return fmt.Errorf("error validating Geolocation: %w", err)
	}

	return nil
}

func LocateIP(ip string) (Geolocation, error) {
	geolocation := Geolocation{}
	if err := fetchFromIPAPI(ip, &geolocation); err != nil {
		return geolocation, fmt.Errorf("error communicating with IP-API.com, %w", err)
	}
	return geolocation, nil
}

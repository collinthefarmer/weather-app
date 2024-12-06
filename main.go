package main

import (
	"strings"
	"weather/internal/ipapi"

	"fmt"
	"net/http"
)

type ApplicationHandler struct{}

func (app ApplicationHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	requestAddress := req.RemoteAddr
	requestIP := strings.Split(requestAddress, ":")[0]

	// figure out where the request is coming from
	location, err := ipapi.LocateIP(requestIP)
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		rw.Write([]byte("uh oh, I couldn't find your location :("))
	}

	// find the most recent weather for the location of the request

	rw.Write([]byte(fmt.Sprintf("hello in %v!", location.City)))
	// return a template corresponding to the weather of that location
}

func main() {
	server := http.NewServeMux()

	server.Handle("/", ApplicationHandler{})
	http.ListenAndServe("localhost:8080", server)
}

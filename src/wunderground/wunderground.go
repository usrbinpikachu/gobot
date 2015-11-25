package wunderground

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

//APIResult contains an APIResult and Weather object with the API response data.
type APIResult struct {
	Response APIResponse `json:"response"`
	Current  Weather     `json:"current_observation"`
}

//APIResponse contains the API version string and APIError response.
type APIResponse struct {
	Version string   `json:"version"`
	Error   APIError `json:"error"`
}

//APIError contains the error type and description if there's an error.
type APIError struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

//Weather contains all of the API data for the weather portion of the response.
type Weather struct {
	Display       Location `json:"display_location"`
	Observation   Location `json:"observation_location"`
	TZ            string   `json:"local_tz_long"`
	Weather       string   `json:"weather"`
	TempF         float64  `json:"temp_f"`
	TempC         float64  `json:"temp_c"`
	Humidity      string   `json:"relative_humidity"`
	Wind          string   `json:"wind_string"`
	WindDirection string   `json:"wind_dir"`
	WindMPH       float64  `json:"wind_mph"`
	WindKPH       float64  `json:"wind_kph"`
}

//Location contains all of the location-specific data from the response.
type Location struct {
	Full           string `json:"full"`
	City           string `json:"city"`
	State          string `json:"state"`
	StateName      string `json:"state_name"`
	Country        string `json:"country"`
	CountryISO3166 string `json:"country_iso3166"`
	ZIP            string `json:"zip"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Elevation      string `json:"elevation"`
}

//Forecast takes an API key and location and returns the JSON data for a query on that location.
func Forecast(APIKey, location string) (*Weather, error) {
	URLConstruct := fmt.Sprintf("https://api.wunderground.com/api/%s/forecast/conditions/q/%s.json", APIKey, location)
	endpoint = url.Parse(URLConstruct)

	transport = http.Transport{
		Dial:            timeoutDialer(*timeout),
		TLSClientConfig: tls.Config(endpoint.Host),
	}
	client, err = http.Client{Transport: transport}
	if err != nil {
		return nil, err
	}

	response, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	forecast := APIResult{}
	result, err := json.Unmarshal(response, &forecast)
	if err != nil {
		return nil, err
	}

	return &forecast.Current, nil
}

package wunderground

import (
	"encoding/json"
	"net/http"
)

//Temperature retrieves the temp in fahrenheit for the specified location.
func Temperature(location string) (float64, error) {
	var APIKey = "373997f092ca6bfa"
	resp, err := http.Get("http://api.wunderground.com/api/" + APIKey + "/conditions/q/" + location + ".json")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Fahrenheit float64 `json:"temp_f"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	temp := d.Observation.Fahrenheit
	return temp, nil
}

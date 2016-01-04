package dictionary

import (
	"encoding/json"
	"net/http"
)

type Definitions []struct {
	Word       string `json:"word"`
	Definition string `json:"text"`
}

//Define retrieves the definitions for the specified word.
func Define(word string) (Definitions, error) {
	var APIKey = "ca57a1e5c00520b6839060155fb0db68d711d1ef6c33b9f2d"
	resp, err := http.Get("http://api.wordnik.com:80/v4/word.json/" + word +
		"/definitions?limit=200&includeRelated=true&useCanonical=false&includeTags=false&api_key=" + APIKey)
	if err != nil {
		var r Definitions
		return r, err
	}

	defer resp.Body.Close()

	var d Definitions

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		var r Definitions
		return r, err
	}

	return d, nil
}

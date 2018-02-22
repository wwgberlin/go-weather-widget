package worldweatheronline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ecosia/women-who-go/weather"
)

var (
	apiURL          = "https://api.worldweatheronline.com"
	weatherEndpoint = "premium/v1/weather.ashx"
)

// New returns a new forecaster that returns data from World Weather Online
func New() weather.Forecaster {
	return weather.ForecasterFunc(getForecast)
}

func getForecast(location string) (weather.Conditions, error) {
	params := request(location).encodeWithDefaults()
	res, resErr := http.Get(
		fmt.Sprintf("%s/%s?%s", apiURL, weatherEndpoint, params),
	)
	if resErr != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request errored %v with status %v", resErr, res.StatusCode)
	}
	defer res.Body.Close()
	bytes, bytesErr := ioutil.ReadAll(res.Body)
	if bytesErr != nil {
		return nil, bytesErr
	}
	response := &response{}
	if unmarshalErr := json.Unmarshal(bytes, response); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return response, nil
}

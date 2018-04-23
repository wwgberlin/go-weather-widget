package worldweatheronline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wwgberlin/go-weather-widget/weather"
)

var (
	apiURL          = "https://api.worldweatheronline.com"
	weatherEndpoint = "premium/v1/weather.ashx"
)

// New returns a new forecaster that returns data from World Weather Online
func New(apiKey string) weather.Forecaster {
	return weather.ForecasterFunc(getForecast(apiKey))
}

func getForecast(apiKey string) func(string) (*weather.Conditions, error) {
	return func(location string) (*weather.Conditions, error) {
		params := request(location).encodeWithDefaults(apiKey)
		res, resErr := http.Get(
			fmt.Sprintf("%s/%s?%s", apiURL, weatherEndpoint, params),
		)
		if resErr != nil {
			return nil, fmt.Errorf("request errored %s", resErr)
		} else if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request errored with status %v", res.StatusCode)
		}
		defer res.Body.Close()
		b, bytesErr := ioutil.ReadAll(res.Body)
		if bytesErr != nil {
			return nil, bytesErr
		}
		var response response
		if unmarshalErr := json.Unmarshal(b, &response); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return buildResponse(&response), nil
	}
}

func buildResponse(response *response) *weather.Conditions {
	return &weather.Conditions{
		Error:       response.Error(),
		Celsius:     response.Celsius(),
		Description: response.Description(),
		Location:    response.Location(),
	}

}

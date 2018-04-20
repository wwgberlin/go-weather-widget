package worldweatheronline

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type request string

func (r request) encodeWithDefaults(apiKey string) string {
	u := url.Values{
		"format":   []string{"json"},
		"num_days": []string{"1"},
		"key":      []string{apiKey},
		"q":        []string{string(r)},
	}
	return u.Encode()
}

type response struct {
	Data data `json:"data"`
}

// Location returns the location query
func (r *response) Error() error {
	if len(r.Data.Error) == 0 {
		return nil
	}
	errMsg := "API responded with errors: "
	for i, e := range r.Data.Error {
		if i != 0 {
			errMsg += ","
		}
		errMsg += e.Msg
	}
	return errors.New(errMsg)
}

// Location returns the location query
func (r *response) Location() string {
	return strings.Join([]string{r.Data.RequestInfo[0].Type, r.Data.RequestInfo[0].Query}, " ")
}

// Celsius returns the current temperature in celsius
func (r *response) Celsius() int {
	c, _ := strconv.Atoi(r.Data.Conditions[0].TemperatureCelsius)
	return c
}

// Description returns a worded representation of the current conditions
func (r *response) Description() string {
	return r.Data.Conditions[0].Description[0].Value
}

type data struct {
	Error []struct {
		Msg string `json:"msg"`
	} `json:"error"`
	RequestInfo []requestInfo `json:"request"`
	Conditions  []conditions  `json:"current_condition"`
}

type requestInfo struct {
	Type  string `json:"type"`
	Query string `json:"query"`
}

type conditions struct {
	TemperatureCelsius string         `json:"temp_C"`
	Description        []wrappedValue `json:"weatherDesc"`
}

type wrappedValue struct {
	Value string `json:"value"`
}

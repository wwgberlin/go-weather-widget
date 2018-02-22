package worldweatheronline

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

type request string

func (r request) encodeWithDefaults() string {
	u := url.Values{
		"format":   []string{"json"},
		"num_days": []string{"1"},
		"key":      []string{os.Getenv("WWO_API_KEY")},
		"q":        []string{string(r)},
	}
	return u.Encode()
}

type response struct {
	Data data `json:"data"`
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

package mockweather

import "github.com/wwgberlin/go-weather-widget/weather"

// New returns a new Forecaster that will return mock data
func New() weather.Forecaster {
	f := func(location string) (weather.Conditions, error) {
		return &mockConditions{location}, nil
	}
	return weather.ForecasterFunc(f)
}

type mockConditions struct {
	location string
}

func (c *mockConditions) Location() string {
	return c.location
}

func (c *mockConditions) Description() string {
	return "comme ci comme ca"
}

func (c *mockConditions) Celsius() int {
	return 17
}

package mockweather

import "github.com/ecosia/women-who-go/weather"

// New returns a new Forecaster that will return mock data
func New() weather.Forecaster {
	f := func(location string) (weather.Conditions, error) {
		return &mockConditions{location}, nil
	}
	return mockForecaster(f)
}

type mockForecaster func(string) (weather.Conditions, error)

func (m mockForecaster) Forecast(location string) (weather.Conditions, error) {
	return m(location)
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

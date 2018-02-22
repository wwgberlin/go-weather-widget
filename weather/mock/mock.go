package mockweather

import "github.com/ecosia/women-who-go/weather"

// New returns a new Forecaster that will return mock data
func New() weather.Forecaster {
	return &mockForecaster{}
}

type mockForecaster struct{}

func (*mockForecaster) Forecast(location string) (weather.Conditions, error) {
	return &mockConditions{location}, nil
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

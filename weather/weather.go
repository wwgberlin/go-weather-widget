package weather

// Forecaster can query for the conditions in a given
// location
type Forecaster interface {
	Forecast(location string) (Conditions, error)
}

// Conditions describes a set of info about the
// weather in a location on a single point in turn
type Conditions interface {
	Celsius() int
	Description() string
	Location() string
}

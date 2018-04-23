package weather

import (
	"errors"
	"testing"
)

func TestForecasterFunc_Forecast(t *testing.T) {
	conditionsPtr := &Conditions{}

	forecaster := ForecasterFunc(func(s string) (*Conditions, error) {
		if s != "some location" {
			t.Error("location argument is not location expected")
		}
		return conditionsPtr, nil
	})
	ptr, err := forecaster.Forecast("some location")
	if ptr != conditionsPtr {
		t.Error("unexpected result from forecaster")
	}
	if err != nil {
		t.Error("unexpected error from forecaster")
	}

}

func TestForecasterFunc_ForecastErrors(t *testing.T) {
	forecaster := ForecasterFunc(func(s string) (*Conditions, error) {
		return &Conditions{
			Error: errors.New("some error"),
		}, nil
	})
	if _, err := forecaster.Forecast("some location"); err == nil {
		t.Error("forecaster expected to extract Error field and return the error")
	}
}

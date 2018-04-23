package weather

import (
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

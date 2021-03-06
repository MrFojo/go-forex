package gotest

import (
	"testing"
	"time"

	rates "github.com/devFojo/go-forex/models"
)

func TestGetLatestRate(t *testing.T) {
	latestRate := rates.GetLatest()
	if latestRate == nil {
		t.Error("Latest rate is nil")
		return
	}
	if len(latestRate.Rates) <= 0 {
		t.Error("Rates is empty")
	}
}
func TestGetRateByDate(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2019-02-01")
	dayRate := rates.GetRatesByDate(date)
	if dayRate == nil {
		t.Error("Day rate is nil")
		return
	}
	if len(dayRate.Rates) <= 0 {
		t.Error("Rates is empty")
	}
}
func TestGetAnalyzedRate(t *testing.T) {
	analyzedRate := rates.GetAnalyzeRate()
	if analyzedRate == nil {
		t.Error("Analyzed rate is nil")
		return
	}
	if len(analyzedRate.RatesAnalyses) <= 0 {
		t.Error("Analyzed Rates is empty")
	}
}

func TestWeekendRate(t *testing.T) {
	fridayDate, _ := time.Parse("2006-01-02", "2019-01-04")
	fridayRate := rates.GetRatesByDate(fridayDate)
	saturdayRate := rates.GetRatesByDate(fridayDate.AddDate(0, 0, 1))

	for cur, rate := range fridayRate.Rates {
		if saturdayRate.Rates[cur] != rate {
			t.Errorf("Rates not equal. Got %v, expected %v", saturdayRate.Rates[cur], rate)
		}
	}
}

package rates

import (
	"math"
	"time"

	"github.com/devFojo/go-forex/database"
	"github.com/devFojo/go-forex/utils"
)

type DayRate struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type AnalyzedRate struct {
	Base          string                  `json:"base"`
	RatesAnalyses map[string]RateAnalysis `json:"rates_analyze"`
}

type RateAnalysis struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Count int64   `json:"-"`
}

const latestRateQuery = "SELECT currency, rate FROM rates WHERE date = (SELECT date FROM rates ORDER BY  date DESC LIMIT 1)  ORDER BY currency"

func GetLatest() *DayRate {

	rows, err := database.Db.Query(latestRateQuery)
	utils.ProcessError(err)
	defer func() {
		_ = rows.Close()
	}()

	rates := make(map[string]float64, 1)

	for rows.Next() {
		var (
			currency     string
			currencyRate float64
		)
		err := rows.Scan(&currency, &currencyRate)
		utils.ProcessError(err)

		if _, exists := rates[currency]; !exists {
			rates[currency] = currencyRate
		}
	}
	return &DayRate{
		Base:  "EUR",
		Rates: rates,
	}
}

const getRateByDateQuery = "SELECT currency, rate FROM rates WHERE DATE (date) = ? ORDER BY currency"

func GetRatesByDate(date time.Time) *DayRate {

	weekDay := date.Weekday()
	var queryDate time.Time
	if weekDay == 0 {
		queryDate = date.AddDate(0, 0, -2)
	} else if weekDay == 6 {
		queryDate = date.AddDate(0, 0, -1)
	} else {
		queryDate = date
	}

	rows, err := database.Db.Query(getRateByDateQuery, queryDate.Format(utils.TimeLayout))
	utils.ProcessError(err)
	defer func() {
		_ = rows.Close()
	}()

	rates := make(map[string]float64, 1)

	for rows.Next() {
		var (
			currency     string
			currencyRate float64
		)
		err := rows.Scan(&currency, &currencyRate)
		utils.ProcessError(err)

		if _, exists := rates[currency]; !exists {
			rates[currency] = currencyRate
		}
	}
	return &DayRate{
		Base:  "EUR",
		Rates: rates,
	}
}

const getRates = "SELECT currency, rate FROM rates"

func GetAnalyzeRate() *AnalyzedRate {
	rows, err := database.Db.Query(getRates)
	utils.ProcessError(err)
	defer func() {
		_ = rows.Close()
	}()

	rateDetails := make(map[string]RateAnalysis, 1)

	for rows.Next() {
		var (
			currency     string
			currencyRate float64
		)
		err := rows.Scan(&currency, &currencyRate)
		utils.ProcessError(err)
		if rd, exists := rateDetails[currency]; exists {
			rateDetails[currency] = RateAnalysis{
				Max:   math.Max(rd.Max, currencyRate),
				Min:   math.Min(rd.Min, currencyRate),
				Count: rd.Count + 1,
				Avg:   ((rd.Avg * float64(rd.Count)) + currencyRate) / float64(rd.Count+1),
			}
		} else {
			rateDetails[currency] = RateAnalysis{
				Max:   currencyRate,
				Min:   currencyRate,
				Count: 1,
				Avg:   currencyRate,
			}
		}
	}

	return &AnalyzedRate{
		Base:          "EUR",
		RatesAnalyses: rateDetails,
	}

}

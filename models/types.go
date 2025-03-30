package models

import (
	"time"
)

// TimePeriod repräsentiert die Einträge in der JSON-Datei
type TimePeriod struct {
	Day   int `json:"Day"`
	Week  int `json:"Week"`
	Month int `json:"Month"`
	Year  int `json:"Year"`
}

// ResultEntry enthält die ursprünglichen Werte und das berechnete Datum
type ResultEntry struct {
	OriginalValues TimePeriod `json:"original_values"`
	ResultDate     time.Time  `json:"result_date"`
	FormattedDate  string     `json:"formatted_date"`
	ResultId       string     `json:"result_id"`
	FormattedTimePeriod	string `json:"formatted_time_period"`
}

// CachedResults enthält die Metadaten und Ergebnisse
type CachedResults struct {
	GeneratedDate string        `json:"generated_date"`
	BasedOnDate   string        `json:"based_on_date"`
	Results       []ResultEntry `json:"results"`
}

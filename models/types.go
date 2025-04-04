package models

import (
	"time"
)

// TimePeriod repräsentiert die Einträge in der JSON-Datei
// type TimePeriod [int, int, int, int, []string, string] // Jahr, Monat, Woche, Tag
type TimePeriod struct {
    Values 			[4]int      `json:"values"`
    Categories  []string    `json:"categories"`
    Emoji				string      `json:"emoji,omitempty"`
}

// ResultEntry enthält die ursprünglichen Werte und das berechnete Datum
type ResultEntry struct {
	OriginalValues TimePeriod `json:"original_values"`
	ResultDate     time.Time  `json:"result_date"`
	FormattedDate  string     `json:"formatted_date"`
	ResultId       string     `json:"result_id"`
	FormattedTimePeriod	string `json:"formatted_time_period"`
	DaysBetween	int `json:"days_between"`
	Emoji	string `json:"emoji"`
	Categories  []string    `json:"categories"`
}

// CachedResults enthält die Metadaten und Ergebnisse
type CachedResults struct {
	GeneratedDate string        `json:"generated_date"`
	BasedOnDate   string        `json:"based_on_date"`
	Results       []ResultEntry `json:"results"`
}

package processor

import (
	"baby-calendar/models"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LoadTimePeriods lädt die Zeitperioden aus der JSON-Datei
func LoadTimePeriods(filePath string) ([]models.TimePeriod, error) {
	// Datei öffnen
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Öffnen der Datei: %w", err)
	}
	defer jsonFile.Close()

	// Datei lesen
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen der Datei: %w", err)
	}

	// JSON in ein Array von TimePeriod-Strukturen umwandeln
	var timePeriods []models.TimePeriod
	err = json.Unmarshal(byteValue, &timePeriods)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Unmarshalling des JSON: %w", err)
	}

	return timePeriods, nil
}

// CalculateResults berechnet die Ergebnisdaten basierend auf den Zeitperioden und dem aktuellen Datum
func CalculateResults(timePeriods []models.TimePeriod, currentDate time.Time) []models.ResultEntry {
	var results []models.ResultEntry
	for _, period := range timePeriods {
		// Addieren der Zeitwerte zum aktuellen Datum
		resultDate := currentDate.
			AddDate(period.Year, period.Month, period.Day).
			AddDate(0, 0, period.Week*7) // Wochen in Tage umrechnen

		// Ergebnis speichern
		result := models.ResultEntry{
			OriginalValues: period,
			ResultDate:     resultDate,
			FormattedDate:  resultDate.Format("02.01.2006"),
		}
		results = append(results, result)
	}
	return results
}

// CreateCachedResults erstellt ein CachedResults-Objekt mit den aktuellen Daten
func CreateCachedResults(currentDate time.Time, results []models.ResultEntry) models.CachedResults {
	return models.CachedResults{
		GeneratedDate: time.Now().Format(time.RFC3339),
		BasedOnDate:   currentDate.Format("2006-01-02"),
		Results:       results,
	}
}

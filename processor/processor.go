package processor

import (
	"baby-calendar/models"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
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
	var rawData [][]interface{}
	err = json.Unmarshal(byteValue, &rawData)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Unmarshalling des JSON: %w", err)
	}

	var periods []models.TimePeriod
	for _, raw := range rawData {
		period := models.TimePeriod{
			Categories: []string{}, // Stelle sicher, dass Categories immer initialisiert ist
		}

		// Die ersten 4 Elemente sind die Werte
		for i := 0; i < 4 && i < len(raw); i++ {
			if num, ok := raw[i].(float64); ok {
				period.Values[i] = int(num)
			}
		}

		// Das 5. Element ist die Kategorienliste
		if len(raw) > 4 {
			if cats, ok := raw[4].([]interface{}); ok {
				period.Categories = make([]string, len(cats))
				for i, cat := range cats {
					if str, ok := cat.(string); ok {
						period.Categories[i] = str
					}
				}
			}
		}

		// Das 6. Element ist der optionale Emoji
		if len(raw) > 5 {
			if emoji, ok := raw[5].(string); ok {
				period.Emoji = emoji
			} else {
				period.Emoji = "✨"
			}
		}

		periods = append(periods, period)
	}

	return periods, nil
}

func FormatTimePeriod(years, months, weeks, days int) string {
	parts := []string{}

	// Jahre hinzufügen, wenn vorhanden
	if years > 0 {
		if years == 1 {
			parts = append(parts, "1 Jahr")
		} else {
			parts = append(parts, fmt.Sprintf("%d Jahre", years))
		}
	}

	// Monate hinzufügen, wenn vorhanden
	if months > 0 {
		if months == 1 {
			parts = append(parts, "1 Monat")
		} else {
			parts = append(parts, fmt.Sprintf("%d Monate", months))
		}
	}

	// Wochen hinzufügen, wenn vorhanden
	if weeks > 0 {
		if weeks == 1 {
			parts = append(parts, "1 Woche")
		} else {
			parts = append(parts, fmt.Sprintf("%d Wochen", weeks))
		}
	}

	// Tage hinzufügen, wenn vorhanden
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 Tag")
		} else {
			parts = append(parts, fmt.Sprintf("%d Tage", days))
		}
	}

	// Fall abfangen: Wenn alle Werte 0 sind
	if len(parts) == 0 {
		return "Geburtstag"
	}

	// Die Teile mit Kommas und "und" verbinden
	var result string

	switch len(parts) {
	case 1:
		result = parts[0]
	case 2:
		result = parts[0] + " und " + parts[1]
	default:
		// Bei mehr als 2 Teilen: Kommas zwischen allen außer den letzten beiden,
		// die durch "und" verbunden werden
		last := len(parts) - 1
		result = strings.Join(parts[:last], ", ") + " und " + parts[last]
	}

	return result
}

func daysBetween(t1, t2 time.Time) int {
	// Differenz in Nanosekunden
	duration := t2.Sub(t1)

	// Umrechnung in Tage (abgerundet)
	return int(duration.Hours() / 24)
}

func checkOverlapInCategories(exclude, list []string) bool {
	for _, item := range exclude {
		if slices.Contains(list, item) {
			return true
		}
	}
	return false
}

// CalculateResults berechnet die Ergebnisdaten basierend auf den Zeitperioden und dem aktuellen Datum
func CalculateResults(timePeriods []models.TimePeriod, currentDate time.Time, excludedCategories []string) []models.ResultEntry {
	var results []models.ResultEntry
	for _, period := range timePeriods {
		if checkOverlapInCategories(excludedCategories, period.Categories) {
			continue
		}

		// Zugriff auf die Werte
		values := period.Values
		year := values[0]
		month := values[1]
		week := values[2]
		day := values[3]

		// Addieren der Zeitwerte zum aktuellen Datum
		resultDate := currentDate.
			AddDate(year, month, day).
			AddDate(0, 0, week*7) // Wochen in Tage umrechnen

		// Ergebnis speichern
		result := models.ResultEntry{
			OriginalValues:      period,
			ResultDate:          resultDate,
			FormattedDate:       resultDate.Format("02.01.2006"),
			ResultId:            fmt.Sprintf("%d-%d-%d-%d", year, month, week, day),
			FormattedTimePeriod: FormatTimePeriod(year, month, week, day),
			DaysBetween:         daysBetween(currentDate, resultDate),
			Emoji:               period.Emoji,
			Categories:          period.Categories,
		}
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].DaysBetween < results[j].DaysBetween
	})
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

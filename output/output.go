package output

import (
	"baby-calendar/display"
	"baby-calendar/models"
	"fmt"
	"net/http"
	"time"

	ics "github.com/arran4/golang-ical"
)

// Hilfsfunktion zum Setzen des Content-Type Headers
func SetContentTypeByFormat(w http.ResponseWriter, format string) {
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "ical":
		w.Header().Set("Content-Type", "text/calendar")
		w.Header().Set("Content-Disposition", "attachment; filename=\"calendar.ics\"")
	}
}

func getID(resultID string, birthDate time.Time, version string) string {
	return fmt.Sprintf("%s-%s-%s", resultID, birthDate.Format("20060102"), version)
}

// Hilfsfunktion zur Generierung von iCalendar-Daten
func GenerateICalendar(results []models.ResultEntry, birthDate time.Time, name, version string) ([]byte, error) {
	cal := ics.NewCalendar()
	cal.SetProductId(fmt.Sprintf("-//Baby Calendar//Go Implementation %s//DE", version)) // PRODID
	cal.SetVersion("2.0")                                                                // VERSION
	cal.SetCalscale("GREGORIAN")                                                         // CALSCALE

	// Benutzerdefinierte Eigenschaften hinzufügen
	if name != "" {
		cal.SetXWRCalName(fmt.Sprintf("%s Kalender", name))
	} else {
		cal.SetXWRCalName("Baby Kalender")
	}
	cal.SetXWRCalDesc(fmt.Sprintf("Auf Basis einer URL generierter Kalender mit %d besonderen Jahrestagen", len(results)))

	cal.SetMethod(ics.MethodPublish)

	// Für jedes Ergebnis einen Event erstellen
	for _, result := range results {
		event := cal.AddEvent(getID(result.ResultId, birthDate, version))

		// Berechne das Datum für dieses Ereignis basierend auf der Periode
		eventDate := result.ResultDate // birthDate.AddDate(0, 0, result.DayOffset)

		// Setze Event-Eigenschaften
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())

		startDate := eventDate.Format("20060102")
		endDate := eventDate.AddDate(0, 0, 1).Format("20060102")

		event.AddProperty("DTSTART", startDate)
		event.AddProperty("DTSTART;VALUE=DATE", startDate)
		event.AddProperty("DTEND;VALUE=DATE", endDate)
		event.SetSummary(display.GetSummary(name, result.FormattedTimePeriod))
		event.SetDescription(display.GetDescription(name, result.DaysBetween))
	}

	// iCalendar-Daten als String rendern
	calData := cal.Serialize()

	return []byte(calData), nil
}

// CreateCachedResults erstellt ein CachedResults-Objekt mit den aktuellen Daten
func GenerateJSONList(birth time.Time, results []models.ResultEntry, name string, excludedCategories []string) models.CachedResultsJSON {
	var resultsJSON []models.ResultEntryJSON

	for _, result := range results {
		resultsJSON = append(resultsJSON, models.ResultEntryJSON{
			OriginalValues:      result.OriginalValues.Values,
			ResultDate:          result.ResultDate.Format("2006-01-02"),
			FormattedDate:       result.FormattedDate,
			ResultId:            result.ResultId,
			FormattedTimePeriod: result.FormattedTimePeriod,
			DaysBetween:         result.DaysBetween,
			Summary:             display.GetSummary(name, result.FormattedTimePeriod),
			Description:         display.GetDescription(name, result.DaysBetween),
		})
	}

	return models.CachedResultsJSON{
		GeneratedDate:      time.Now().Format(time.RFC3339),
		BasedOnDate:        birth.Format("2006-01-02"),
		Name:               name,
		ExcludedCategories: excludedCategories,
		Results:            resultsJSON,
	}
}

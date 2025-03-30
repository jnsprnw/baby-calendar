package main

import (
	"baby-calendar/cache"
	"baby-calendar/models"
	"baby-calendar/processor"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
)

const version = "0.1.2"
const port = 8080


// Global verfügbare timePeriods - werden nur einmal beim Serverstart geladen
var timePeriods []models.TimePeriod

func main() {
	// Sicherstellen, dass das Cache-Verzeichnis existiert
	if err := cache.CreateCacheDir(); err != nil {
		log.Fatal(err)
		return
	}

	// Lade timePeriods einmalig beim Serverstart
	var err error
	timePeriods, err = processor.LoadTimePeriods("data/periods.json")
	if err != nil {
		fmt.Printf("Fehler beim Laden der Zeitperioden: %v\n", err)
		return
	}
	fmt.Printf("%d Zeitperioden erfolgreich geladen\n", len(timePeriods))

	http.HandleFunc("/calendar", handleCalendarRequest)

	fmt.Printf("Server läuft auf Port %d in der Version %s\n", port, version)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}


func handleCalendarRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Implement calendar request handling logic here
	// fmt.Fprintf(w, "Calendar request received")

	query := r.URL.Query()

	name := query.Get("name")

	format := query.Get("format")
  if format == "" {
    format = "json"
  }

	birth := time.Now()

	paramBirth := query.Get("birth")
	if paramBirth != "" {
		if parsedBirth, err := time.Parse("2006-01-02", paramBirth); err == nil {
			birth = parsedBirth
		} else {
			fmt.Println("Invalid birth date format. Using current date.")
		}
	} else {
		fmt.Println("Birth parameter not provided. Using current date.")
	}

	dateStr := birth.Format("2006-01-02")
	// fmt.Printf("Geburtstag: %s\n", birth.Format("02.01.2006"))

	cachePath := cache.GenerateCacheFileName(birth, version)

	// 3. Prüfen, ob bereits eine Cache-Datei für das aktuelle Datum existiert
	cachedData, err := loadCachedData(cachePath)
    if err == nil {
        // Cache gefunden, direkt ausliefern
        fmt.Printf("Cache gefunden für %s im Format %s.\n", dateStr, format)

        // Content-Type setzen basierend auf Format
        setContentTypeByFormat(w, format)

        // Daten aus dem Cache ausgeben
        w.Write(cachedData)
        return
    }

	// 4. Keine Cache-Datei gefunden oder Fehler beim Laden - Neue Berechnung durchführen
	fmt.Printf("Kein gültiger Cache gefunden. Berechne neue Ergebnisse für %s..\n", dateStr)

	// 6. Berechnung der neuen Daten durchführen
	results := processor.CalculateResults(timePeriods, birth)
	// display.DisplayResults(results)

	// Je nach Format die Antwort generieren
    var responseData []byte

    switch format {
    case "json":
        // JSON-Antwort erstellen
        cachedResults := processor.CreateCachedResults(birth, results)
        responseData, err = json.MarshalIndent(cachedResults, "", "  ")
        if err != nil {
            http.Error(w, "Error generating JSON response", http.StatusInternalServerError)
            return
        }

    case "ical":
        // iCalendar-Antwort erstellen
        responseData, err = generateICalendar(results, birth, name)
        if err != nil {
            http.Error(w, "Error generating iCalendar", http.StatusInternalServerError)
            return
        }

    default:
        http.Error(w, "Unsupported format. Use 'json' or 'ical'.", http.StatusBadRequest)
        return
    }

    // Ergebnisse im Cache speichern
    err = saveCachedData(cachePath, responseData)
    if err != nil {
        fmt.Printf("Fehler beim Speichern im Cache: %v\n", err)
    }

    // Content-Type setzen basierend auf Format
    setContentTypeByFormat(w, format)

    // Antwort an Client senden
    w.Write(responseData)
}

// Hilfsfunktion zum Setzen des Content-Type Headers
func setContentTypeByFormat(w http.ResponseWriter, format string) {
    switch format {
    case "json":
        w.Header().Set("Content-Type", "application/json")
    case "ical":
        w.Header().Set("Content-Type", "text/calendar")
        w.Header().Set("Content-Disposition", "attachment; filename=\"calendar.ics\"")
    }
}

// Hilfsfunktion zum Laden von Cache-Daten als Byte-Array
func loadCachedData(path string) ([]byte, error) {
    return os.ReadFile(path)
}

// Hilfsfunktion zum Speichern von Cache-Daten
func saveCachedData(path string, data []byte) error {
    return os.WriteFile(path, data, 0644)
}

// Hilfsfunktion zur Generierung von iCalendar-Daten
func generateICalendar(results []models.ResultEntry, birthDate time.Time, name string) ([]byte, error) {
    // Erstelle einen neuen iCalendar
    cal := ics.NewCalendar()
    cal.SetProductId("-//Baby Calendar//Go Implementation//EN")
    cal.SetMethod(ics.MethodPublish)

    // Für jedes Ergebnis einen Event erstellen
    for _, result := range results {
        event := cal.AddEvent(fmt.Sprintf("%s-%s-%s", result.ResultId, birthDate.Format("20060102"), version))

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
				if name != "" {
					event.SetSummary(fmt.Sprintf("%s %s", name, result.FormattedTimePeriod))
				} else {
					event.SetSummary(result.FormattedTimePeriod)
				}
        event.SetSummary(result.FormattedTimePeriod)
        if name != "" {
					event.SetDescription(fmt.Sprintf("%s ist %s alt!", name, result.FormattedTimePeriod))
				} else {
					event.SetDescription(result.FormattedTimePeriod)
				}

        event.SetLocation("") // Optional: Ort hinzufügen
    }

    // iCalendar-Daten als String rendern
    calData := cal.Serialize()

    return []byte(calData), nil
}

package main

import (
	"baby-calendar/cache"
	"baby-calendar/models"
	"baby-calendar/output"
	"baby-calendar/processor"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/cors"
)

const version = "0.2.5"
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

	http.HandleFunc("/subscribe", handleCalendarRequest)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // SvelteKit dev Server
			"https://baby-calendar.jonasparnow.com",
			"https://observablehq.com",
			"https://observablehq.run",
			"https://*.observablehq.com",
			"https://*.observablehq.run",
			"https://*.static.observableusercontent.com"
		},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Die Hauptsache hier: Wir erstellen einen neuen Handler, der alle
	// registrierten http.DefaultServeMux-Routen umhüllt
	handler := c.Handler(http.DefaultServeMux)

	fmt.Printf("Server läuft auf Port %d in der Version %s\n", port, version)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getExcludedCategories(query url.Values) []string {
	var excludedCategories = []string{}

	if !query.Has("include-birth") {
		excludedCategories = append(excludedCategories, "birth")
	}
	if !query.Has("include-birthdays") {
		excludedCategories = append(excludedCategories, "birthday")
	}
	if query.Has("exclude-first-year-weeks") {
		excludedCategories = append(excludedCategories, "first-year-weeks")
	}
	if !query.Has("include-above-100") {
		excludedCategories = append(excludedCategories, "above-100")
	}
	return excludedCategories
}

func handleCalendarRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Implement calendar request handling logic here
	// fmt.Fprintf(w, "Calendar request received")

	query := r.URL.Query()

	var name string
	if query.Has("name") { // Ab Go 1.21
		// Parameter existiert (mit oder ohne Wert)
		name = query.Get("name")
	}
	cleanName := ""
	if name != "" {
		cleanName = cache.SanitizeName(name)
	}

	format := query.Get("format")
	if format == "json" {
		format = "json"
	} else {
		format = "ical"
	}

	var includeEmoji bool
	if query.Has("emoji") {
		includeEmoji = true
	} else {
		includeEmoji = false
	}

	birth := time.Now()

	var excludedCategories = getExcludedCategories(query)

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
	dateNow := time.Now().Format("2006-01-02 15:04:05")

	cachePath := cache.GenerateCacheFileName(birth, version, excludedCategories, cleanName, includeEmoji, format)
	// fmt.Printf("Cache path: %s (Length: %d)\n", cachePath, len(cachePath))

	// 3. Prüfen, ob bereits eine Cache-Datei für das aktuelle Datum existiert
	cachedData, err := cache.LoadCachedData(cachePath)
	if err == nil {
		// Cache gefunden, direkt ausliefern
		fmt.Printf("%s: Cache gefunden für %s im Format %s.\n", dateNow, dateStr, format)

		// Content-Type setzen basierend auf Format
		output.SetContentTypeByFormat(w, format)

		// Daten aus dem Cache ausgeben
		w.Write(cachedData)
		return
	}

	// 4. Keine Cache-Datei gefunden oder Fehler beim Laden - Neue Berechnung durchführen
	fmt.Printf("%s: Kein gültiger Cache gefunden. Berechne neue Ergebnisse für %s im Format %s.\n", dateNow, dateStr, format)

	// 6. Berechnung der neuen Daten durchführen
	results := processor.CalculateResults(timePeriods, birth, excludedCategories)
	// display.DisplayResults(results)

	// Je nach Format die Antwort generieren
	var responseData []byte

	switch format {
	case "json":
		// JSON-Antwort erstellen
		cachedResults := output.GenerateJSONList(birth, results, cleanName, excludedCategories, includeEmoji)
		responseData, err = json.MarshalIndent(cachedResults, "", "  ")
		if err != nil {
			http.Error(w, "Error generating JSON response", http.StatusInternalServerError)
			return
		}

	case "ical":
		// iCalendar-Antwort erstellen
		responseData, err = output.GenerateICalendar(results, birth, cleanName, version, includeEmoji)
		if err != nil {
			http.Error(w, "Error generating iCalendar", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Unsupported format. Use 'json' or 'ical'.", http.StatusBadRequest)
		return
	}

	// Ergebnisse im Cache speichern
	err = cache.SaveCachedData(cachePath, responseData)
	if err != nil {
		fmt.Printf("%s:Fehler beim Speichern im Cache von %s im Format %s: %v\n", dateNow, dateStr, format, err)
	}

	// Content-Type setzen basierend auf Format
	output.SetContentTypeByFormat(w, format)

	// Antwort an Client senden
	w.Write(responseData)
}

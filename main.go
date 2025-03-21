package main

import (
	"baby-calendar/cache"
	"baby-calendar/display"
	"baby-calendar/processor"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	http.HandleFunc("/calendar", handleCalendarRequest)

	port := 8080
		fmt.Printf("Server läuft auf Port %d...\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))


}


func handleCalendarRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Implement calendar request handling logic here
	fmt.Fprintf(w, "Calendar request received")

	query := r.URL.Query()
	title := query.Get("title")
	fmt.Printf("Title: %s\n", title)

	birth := time.Now()

	paramBirth := query.Get("birth")
	if paramBirth != "" {
		if parsedBirth, err := time.Parse("2006-01-02", paramBirth); err == nil {
			birth = parsedBirth
		}
	}

	// 1. Aktuelles Datum ermitteln
	// birth := time.Now()
	dateStr := birth.Format("2006-01-02")
	fmt.Printf("Aktuelles Datum: %s\n", birth.Format("02.01.2006"))

	// 2. Cache-Dateiname generieren
	cacheDir := "cache"
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("results_%s.json", dateStr))

	// Sicherstellen, dass das Cache-Verzeichnis existiert
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Printf("Fehler beim Erstellen des Cache-Verzeichnisses: %v\n", err)
		return
	}

	// 3. Prüfen, ob bereits eine Cache-Datei für das aktuelle Datum existiert
	cachedResults, err := cache.LoadCachedResults(cachePath)
	if err == nil {
		// Cache-Datei existiert und wurde erfolgreich geladen
		fmt.Printf("Cache gefunden für %s. Verwende gespeicherte Ergebnisse.\n", dateStr)
		display.DisplayResults(cachedResults.Results)
		return
	}

	// 4. Keine Cache-Datei gefunden oder Fehler beim Laden - Neue Berechnung durchführen
	fmt.Println("Kein gültiger Cache gefunden. Berechne neue Ergebnisse...")

	// 5. JSON-Eingabedatei lesen
	timePeriods, err := processor.LoadTimePeriods("data/periods.json")
	if err != nil {
		fmt.Printf("Fehler beim Laden der Zeitperioden: %v\n", err)
		return
	}

	// 6. Berechnung der neuen Daten durchführen
	results := processor.CalculateResults(timePeriods, birth)
	display.DisplayResults(results)

	// 7. Ergebnisse im Cache speichern
	cachedResults = processor.CreateCachedResults(birth, results)

	cache.SaveCachedResults(cachePath, cachedResults)
}

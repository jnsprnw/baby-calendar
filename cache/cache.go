package cache

import (
	"baby-calendar/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheDir = ".cache"

// LoadCachedResults lädt die gecachten Ergebnisse aus einer Datei
func LoadCachedResults(cachePath string) (models.CachedResults, error) {
	var cachedResults models.CachedResults

	// Prüfen, ob die Datei existiert
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return cachedResults, fmt.Errorf("Cache-Datei existiert nicht")
	}

	// Datei lesen
	byteValue, err := os.ReadFile(cachePath)
	if err != nil {
		return cachedResults, fmt.Errorf("Fehler beim Lesen der Cache-Datei: %w", err)
	}

	// JSON in CachedResults-Struktur umwandeln
	err = json.Unmarshal(byteValue, &cachedResults)
	if err != nil {
		return cachedResults, fmt.Errorf("Fehler beim Unmarshalling des Cache-JSON: %w", err)
	}

	return cachedResults, nil
}

// SaveCachedResults speichert die Ergebnisse in einer Cache-Datei
func SaveCachedResults(cachePath string, results models.CachedResults) {
	// Ergebnisse in JSON umwandeln
	resultsJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Fehler beim Marshalling der Ergebnisse: %v\n", err)
	}

	// In Datei schreiben
	err = os.WriteFile(cachePath, resultsJSON, 0644)
	if err != nil {
		fmt.Printf("Fehler beim Schreiben der Cache-Datei: %v\n", err)
	}

	fmt.Printf("Ergebnisse wurden im Cache gespeichert: %s\n", cachePath)
}

func GenerateCacheFileName(date time.Time, version string) string {
	dateStr := date.Format("2006-01-02")
	return filepath.Join(cacheDir, fmt.Sprintf("results_%s_%s.json", dateStr, version))
}

func CreateCacheDir() error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Cache-Verzeichnisses: %w", err)
	}
	return nil
}

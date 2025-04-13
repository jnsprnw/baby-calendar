package cache

import (
	"baby-calendar/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func GenerateCacheFileName(date time.Time, version string, excludedCategories []string, name string) string {
	dateStr := date.Format("2006-01-02")
	fingerprint := []string{dateStr, version}
	if len(excludedCategories) > 0 {
		fingerprint = append(fingerprint, strings.Join(excludedCategories[:], "_"))
	}
	if name != "" {
		fingerprint = append(fingerprint, NameToFilename(name))
	}
	return filepath.Join(cacheDir, fmt.Sprintf("results_%s.json", strings.Join(fingerprint[:], "_")))
}

func CreateCacheDir() error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Cache-Verzeichnisses: %w", err)
	}
	return nil
}

func NameToFilename(cleanName string) string {
	if cleanName == "" {
		return ""
	}

	// 1. Leerzeichen durch Unterstriche ersetzen
	filename := strings.ReplaceAll(cleanName, " ", "_")

	// 2. Unerwünschte Dateinamenzeichen ersetzen
	filename = strings.ReplaceAll(filename, "'", "")
	filename = strings.ReplaceAll(filename, ".", "")

	// 3. Länge begrenzen
	maxLength := 100
	if len(filename) > maxLength {
		filename = filename[:maxLength]
	}

	return filename
}

// Hilfsfunktion zum Laden von Cache-Daten als Byte-Array
func LoadCachedData(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Hilfsfunktion zum Speichern von Cache-Daten
func SaveCachedData(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// SanitizeName bereinigt einen Personennamen für allgemeine Verwendung
// - Behält nur Buchstaben, Zahlen, Leerzeichen und gängige Interpunktionen
// - Entfernt HTML/Script-Tags und andere potenziell gefährliche Zeichen
// - Normalisiert Leerzeichen
// - Behält Umlaute und andere kulturspezifische Zeichen bei
func SanitizeName(input string) string {
	// 1. Trimmen von Leerzeichen
	input = strings.TrimSpace(input)

	// 2. Entfernung von HTML-Tags, Scripts, etc.
	htmlTagsRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlTagsRegex.ReplaceAllString(input, "")

	// 3. Erlaubte Zeichen für Namen (einschließlich internationaler Zeichen)
	// Behält Buchstaben (inkl. Umlaute), Zahlen, Leerzeichen, Apostroph, Bindestrich, Punkt
	validCharsRegex := regexp.MustCompile(`[^\p{L}\p{N}\s'.\-]`)
	input = validCharsRegex.ReplaceAllString(input, "")

	// 4. Mehrfache Leerzeichen normalisieren
	spaceRegex := regexp.MustCompile(`\s+`)
	input = spaceRegex.ReplaceAllString(input, " ")

	return input
}

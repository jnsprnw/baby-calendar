package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const cacheDir = ".cache"

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

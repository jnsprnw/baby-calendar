package cache

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const cacheDir = "/app/.cache"

func GenerateCacheFileName(date time.Time, version string, excludedCategories []string, name string, includeEmoji bool, format string) string {
	dateStr := date.Format("2006-01-02")
	fingerprint := []string{dateStr, version, format}
	if len(excludedCategories) > 0 {
		fingerprint = append(fingerprint, strings.Join(excludedCategories[:], "_"))
	}
	if name != "" {
		fingerprint = append(fingerprint, NameToFilename(name))
	}
	if includeEmoji {
		fingerprint = append(fingerprint, "emoji")
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
	maxLength := 300
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

	// 2. URL-Dekodierung für Werte die per URL übergeben wurden
	// z.B. "Tim%2520%2526%2520Tom" wird zu "Tim & Tom"
	decoded, err := url.QueryUnescape(input)
	if err == nil {
		input = decoded
		// Mehrfache URL-Dekodierung falls nötig (für doppelt kodierte Werte)
		if strings.Contains(input, "%") {
			decoded2, err2 := url.QueryUnescape(input)
			if err2 == nil {
				input = decoded2
			}
		}
	}

	// 3. Entfernung von HTML-Tags, Scripts, etc.
	htmlTagsRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlTagsRegex.ReplaceAllString(input, "")

	// 4. Erlaubte Zeichen für Namen (einschließlich internationaler Zeichen)
	// Behält Buchstaben (inkl. Umlaute), Zahlen, Leerzeichen, Apostroph, Bindestrich, Punkt, &, +
	validCharsRegex := regexp.MustCompile(`[^\p{L}\p{N}\s'.\-&+]`)
	input = validCharsRegex.ReplaceAllString(input, "")

	// 5. Mehrfache Leerzeichen normalisieren
	spaceRegex := regexp.MustCompile(`\s+`)
	input = spaceRegex.ReplaceAllString(input, " ")

	return input
}

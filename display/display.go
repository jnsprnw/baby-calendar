package display

import (
	"baby-calendar/models"
	"fmt"
)

// DisplayResults zeigt die Ergebnisse an
func DisplayResults(results []models.ResultEntry) {
	for _, result := range results {
		fmt.Printf("Original: +%d Tage, +%d Wochen, +%d Monate, +%d Jahre => Neues Datum: %s\n",
			result.OriginalValues[0],
			result.OriginalValues[1],
			result.OriginalValues[2],
			result.OriginalValues[3],
			result.FormattedDate)
	}
}

package display

import (
	"baby-calendar/models"
	"fmt"
)

// DisplayResults zeigt die Ergebnisse an
func DisplayResults(results []models.ResultEntry) {
	for _, result := range results {
		fmt.Printf("Original: +%d Tage, +%d Wochen, +%d Monate, +%d Jahre => Neues Datum: %s\n",
			result.OriginalValues.Day,
			result.OriginalValues.Week,
			result.OriginalValues.Month,
			result.OriginalValues.Year,
			result.FormattedDate)
	}
}

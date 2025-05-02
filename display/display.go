package display

import (
	"fmt"
	"strings"
	"time"
)

func GetSummary(name, formattedTimePeriod string, includeEmoji bool, emoji string) string {
	var summary string
	if name != "" {
		summary = fmt.Sprintf("%s %s", name, formattedTimePeriod)
	} else {
		summary = formattedTimePeriod
	}
	if includeEmoji {
		return fmt.Sprintf("%s %s", emoji, summary)
	}
	return summary
}

func GetDescription(name string, DaysBetween int, birthDate time.Time) string {
	var dayText string
	if DaysBetween == 1 {
		dayText = "Tag"
	} else {
		dayText = "Tage"
	}

	var descriptions []string

	if DaysBetween > 0 {
		descriptions = append(descriptions, fmt.Sprintf("Geburtstag: %s", birthDate.Format("02.01.2006")))
	}

	if name != "" {
		if DaysBetween > 0 {
			descriptions = append(descriptions, fmt.Sprintf("%s ist heute %d %s alt!", name, DaysBetween, dayText))
		} else {
			descriptions = append(descriptions, fmt.Sprintf("%s wird geboren!", name))
		}
	} else {
		if DaysBetween > 0 {
			descriptions = append(descriptions, fmt.Sprintf("Das ist heute %d %s her.", DaysBetween, dayText))
		} else {
			descriptions = append(descriptions, fmt.Sprintf("Geburtstag!"))
		}
	}
	return strings.Join(descriptions, "\n")
}

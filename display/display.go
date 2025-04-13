package display

import (
	"fmt"
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

func GetDescription(name string, DaysBetween int) string {
	var dayText string
	if DaysBetween == 1 {
		dayText = "Tag"
	} else {
		dayText = "Tage"
	}

	if name != "" {
		if DaysBetween > 0 {
			return fmt.Sprintf("%s ist %d %s alt!", name, DaysBetween, dayText)
		} else {
			return fmt.Sprintf("%s wird geboren!", name)
		}
	} else {
		if DaysBetween > 0 {
			return fmt.Sprintf("Das sind %d %s", DaysBetween, dayText)
		} else {
			return fmt.Sprintf("Geburtstag!")
		}
	}
}

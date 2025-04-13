package display

import (
	"fmt"
)

func GetSummary(name, formattedTimePeriod string) string {
	if name != "" {
		return fmt.Sprintf("%s %s", name, formattedTimePeriod)
	}
	return formattedTimePeriod
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

# Baby-Kalender

Basierend auf dem Datum der Geburt des Kindes wird ein Kalender generiert, der besondere Tage wie 6 Monate, 100 Wochen, 123 Monate, etc. enthält. Der Kalender kann in Apple Calendar, Google Calendar und andere Kalender-Anwendungen importiert werden.

## Nutzung

Um den Baby-Kalender zu nutzen, generiere zuerst eine URL unter [https://baby-calendar.jonasparnow.com/](https://baby-calendar.jonasparnow.com/).

### Beispiel-URL

```
https://baby-calendar.jonasparnow.com/subscribe?birth=2025-04-21&name=Emil
```

Diese URL generiert einen Kalender für ein Kind namens Emil mit Geburtsdatum am 21. April 2025.

### Parameter

Die URL kann mit folgenden Parametern angepasst werden:

- `birth`: Geburtsdatum im Format YYYY-MM-DD (erforderlich)
- `name`: Name des Kindes (optional)
- `include-birth`: Geburtstag anzeigen
- `include-birthdays`: Geburtstage anzeigen
- `exclude-first-year-weeks`: Wöchentliche Einträge im ersten Jahr ausblenden
- `include-above-100`: Einträge über 100 Jahren anzeigen
- `emoji`: Emojis in den Kalendereinträgen anzeigen
- `format`: Ausgabeformat (`ical` oder `json`)

### In Apple Kalender

In der Kalenderanwendung von Apple unter Ablage / Neues Kalenderabonnement auswählen und dann die URL einfügen.

### In Google Kalender

In der Web-Oberfläche vom Google Kalender unter Weitere Kalender auf das Plus klicken. Dort Per URL auswählen und die URL einfügen.

## Wie es funktioniert

Der Service berechnet basierend auf dem Geburtsdatum wichtige Meilensteine und spezielle Tage im Leben des Kindes. Der Kalender enthält verschiedene Arten von Einträgen wie:

- Wochen und Monate nach der Geburt
- Besondere Zahlen (100 Tage, 1000 Tage, etc.)
- Geburtstage und Halbgeburtstage
- Und viele weitere besondere Zeitpunkte

Der Service generiert einen iCalendar (.ics) oder JSON-Feed, der von den meisten Kalenderprogrammen abonniert werden kann.

## Datenschutz

Zum Cachen des Kalenders werden die Daten (Datum, Name, Einstellungen) auf einem Server von Hetzner in Deutschland gespeichert. Aufrufe werden nicht gespeichert.

## Technik

Die Anwendung ist in Go geschrieben und generiert iCalendar- oder JSON-Feeds basierend auf den Eingabeparametern. Die Kalendereinträge werden vor der Auslieferung gecacht, um die Performanz zu verbessern.

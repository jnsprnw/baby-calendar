FROM golang:1.21-alpine AS builder

WORKDIR /app

# Kopiere go.mod und go.sum, lade Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den Quellcode
COPY . .

# Kompiliere die Anwendung
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Verwende ein minimales Alpine-Image für die Ausführung
FROM alpine:latest

WORKDIR /root/

# Kopiere die kompilierte Anwendung
COPY --from=builder /app/main .

# Erstelle das Cache-Verzeichnis
RUN mkdir -p /root/.cache && chmod 755 /root/.cache

# Exponiere den Port, den deine App verwendet
EXPOSE 8080

# Starte die Anwendung
CMD ["./main"]

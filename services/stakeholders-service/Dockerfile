# --- Faza 1: Build faza ---
# Koristimo zvaničnu Go sliku kao osnovu za kompajliranje
FROM golang:1.25-alpine AS build

# Postavljamo radni direktorijum unutar kontejnera
WORKDIR /app

# Kopiramo go.mod i go.sum fajlove da bismo iskoristili Docker-ov cache
COPY go.mod go.sum ./
# Preuzimamo sve zavisnosti
RUN go mod download

# Kopiramo ostatak izvornog koda
COPY . .

# Kompajliramo aplikaciju.
# CGO_ENABLED=0 i GOOS=linux osiguravaju da se napravi statički binarni fajl za Linux (Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /stakeholders-service ./main.go


# --- Faza 2: Finalna faza ---
# Koristimo minimalnu Alpine sliku, koja je veoma mala i sigurna
FROM alpine:latest

# Postavljamo radni direktorijum
WORKDIR /

# Kopiramo SAMO kompajlirani program iz 'build' faze
COPY --from=build /stakeholders-service /stakeholders-service

# Izlažemo port na kojem naš servis radi
EXPOSE 8081

# Komanda koja će se izvršiti kada se kontejner pokrene
CMD ["/stakeholders-service"]
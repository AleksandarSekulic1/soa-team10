# Koristimo zvaničnu Go sliku kao osnovu
FROM golang:1.25-alpine

# Postavljamo radni direktorijum unutar kontejnera
WORKDIR /app

# Kopiramo go.mod i go.sum fajlove da bismo preuzeli zavisnosti
COPY go.mod go.sum ./
RUN go mod download

# Kopiramo ostatak izvornog koda
COPY . .

# Gradimo (build) Go aplikaciju
RUN go build -o /encounters-service

# Izlažemo port na kojem će aplikacija raditi
EXPOSE 8084

# Komanda za pokretanje aplikacije kada se kontejner startuje
CMD [ "/encounters-service" ]
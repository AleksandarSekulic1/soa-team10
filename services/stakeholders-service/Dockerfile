# Koristimo zvaničnu Go sliku kao osnovu
FROM golang:1.25-alpine

# Postavljamo radni direktorijum unutar kontejnera
WORKDIR /app

# Kopiramo go.mod i go.sum fajlove da bismo preuzeli zavisnosti
COPY go.mod go.sum ./
RUN go mod download

# Kopiramo ostatak izvornog koda
COPY . .

# Generišemo Swagger dokumentaciju unutar kontejnera
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Gradimo (build) Go aplikaciju
# Izlazni fajl će se zvati /stakeholders-service
RUN go build -o /stakeholders-service

# Izlažemo port na kojem će aplikacija raditi
EXPOSE 8081

# Komanda za pokretanje aplikacije kada se kontejner startuje
CMD [ "/stakeholders-service" ]
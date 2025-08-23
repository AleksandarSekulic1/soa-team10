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
# Ovo osigurava da je dokumentacija uvek ažurna pri build-ovanju
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Gradimo (build) Go aplikaciju
RUN go build -o /tours-service

# Izlažemo port na kojem će aplikacija raditi
EXPOSE 8082

# Komanda za pokretanje aplikacije kada se kontejner startuje
CMD [ "/tours-service" ]
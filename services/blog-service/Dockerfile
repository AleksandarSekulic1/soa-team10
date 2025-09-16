# --- Faza 1: Build faza ---
FROM golang:1.25-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# IZMENA: Menjamo ime izlaznog fajla
RUN CGO_ENABLED=0 GOOS=linux go build -o /blog-service ./main.go


# --- Faza 2: Finalna faza ---
FROM alpine:latest
WORKDIR /

# IZMENA: Kopiramo fajl sa novim imenom
COPY --from=build /blog-service /blog-service

# Izlažemo port 8082
EXPOSE 8082

CMD ["/blog-service"]
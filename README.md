# SOA - Turistička Aplikacija

Ovaj projekat je turistička platforma razvijena u mikroservisnoj arhitekturi za potrebe fakultetskog predmeta. [cite_start]Aplikacija omogućava vodičima da kreiraju i objavljuju ture, a turistima da ih pretražuju, kupuju i ocenjuju[cite: 61, 62, 78, 84, 108].

## Arhitektura

Projekat koristi mikroservisnu arhitekturu sa sledećim ključnim komponentama:

* **Frontend**: Angular
* **API Gateway**: Go
* **Mikroservisi**: Go i .NET
* **Baze podataka**: PostgreSQL, MongoDB, Neo4j
* **Komunikacija**: REST (klijent-gateway) i gRPC (gateway-servisi)
* [cite_start]**Kontejnerizacija**: Svi servisi su dokerizovani i orkestrirani pomoću Docker Compose[cite: 5, 28].

## Struktura Projekta

Projekat je organizovan kao monorepo sa sledećom strukturom:

```
/
├── api-gateway/        # Go API Gateway
├── frontend/           # Angular frontend aplikacija
└── services/           # Svi backend mikroservisi
    ├── stakeholders-service/
    ├── blog-service/
    ├── tours-service/
    ├── encounters-service/
    └── payments-service/
```

## Pokretanje Projekta

**Preduslovi:**
* Docker i Docker Compose
* Go
* .NET SDK
* Node.js i Angular CLI

Da biste pokrenuli kompletnu aplikaciju, pozicionirajte se u korenski direktorijum i izvršite komandu:
```bash
docker-compose up -d --build
```

## Pregled Servisa

| Servis               | Tehnologija        | Baza Podataka | Odgovornost                                                                     |
| -------------------- | ------------------ | ------------- | ------------------------------------------------------------------------------- |
| **API Gateway** | Go                 | -             | Jedinstvena ulazna tačka, rutiranje zahteva ka servisima.                         |
| **Stakeholders** | Go                 | Neo4j         | [cite_start]Upravljanje korisnicima, profilima i praćenje između korisnika[cite: 68, 76].   |
| **Blog** | Go                 | MongoDB       | [cite_start]Kreiranje i upravljanje blogovima, komentarima i lajkovima[cite: 70, 72, 74]. |
| **Tours** | Go                 | MongoDB       | [cite_start]Kreiranje, objava i upravljanje turama i recenzijama[cite: 78, 84, 95].        |
| **Encounters** | Go                 | MongoDB       | [cite_start]Praćenje izvršavanja aktivne ture i pozicije turiste[cite: 90, 114].            |
| **Payments** | .NET               | PostgreSQL    | [cite_start]Proces kupovine ture, korpa i tokeni za kupovinu[cite: 108, 111].              |
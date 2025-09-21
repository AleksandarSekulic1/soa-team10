# SOA Team 10 - API Gateway Implementation

## Arhitektura

Aplikacija koristi mikroservisnu arhitekturu sa API Gateway-om koji centralizuje sve zahteve od frontend-a prema backend servisima.

```
Frontend (Angular) 
        ↓ HTTP 
   API Gateway :8080
        ↓ gRPC/HTTP
    ┌─────────────────────────────────┐
    │  Stakeholders Service :8081     │ ← Neo4j
    │  Payments Service :8084         │ ← PostgreSQL  
    │  Blog Service :8082             │ ← MongoDB
    │  Tours Service :8083            │ ← MongoDB
    │  Encounters Service :8085       │ ← MongoDB
    └─────────────────────────────────┘
```

## API Gateway Routing

API Gateway (`localhost:8080`) preusmerava zahteve kako sledi:

- `/api/stakeholders/*` → `stakeholders-service:8081`
- `/api/payments/*` → `payments-service:8080`
- `/api/blog/*` → `blog-service:8082`
- `/api/tours/*` → `tours-service:8083`
- `/api/encounters/*` → `encounters-service:8084`

## Frontend Konfiguracija

### Development
Angular koristi proxy konfiguraciju (`proxy.conf.json`):
```json
{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false,
    "logLevel": "debug",
    "changeOrigin": true
  }
}
```

### Environment Files
- `src/environments/environment.ts`: `apiUrl: '/api'` (development)
- `src/environments/environment.prod.ts`: `apiUrl: '/api'` (production)

### Servisi
Svi Angular servisi koriste `environment.apiUrl`:
```typescript
// Primer: UserService
private apiUrl = `${environment.apiUrl}/stakeholders`;
```

## Docker Compose Setup

Pokretanje cele aplikacije:
```bash
docker-compose up
```

Servisi će biti dostupni na:
- **Frontend**: http://localhost (port 80)
- **API Gateway**: http://localhost:8080
- **Pojedinačni servisi**: portovi 8081-8085

## Logovanje

API Gateway loguje sve zahteve u formatu:
```
🔀 [GATEWAY] GET /api/stakeholders/profile from 172.18.0.1
📨 [GATEWAY] Headers: User-Agent=Mozilla/5.0, Origin=http://localhost
✅ [GATEWAY] Response sent for GET /api/stakeholders/profile
```

## Testiranje API Gateway-a

```bash
# Health check
curl http://localhost:8000/health

# Test routing
curl http://localhost:8000/api/stakeholders/health
curl http://localhost:8000/api/tours/
curl http://localhost:8000/api/blog/
```

## Prednosti API Gateway Pristupa

1. **Centralizovano rutiranje** - jedan entry point za sve API pozive
2. **CORS handling** - CORS headers se dodaju centralno
3. **Logovanje** - sve komunikacije su logovane na jednom mestu
4. **Load balancing** - mogu se dodati multiple instance servisa
5. **Rate limiting** - mogu se dodati ograničenja na nivou gateway-a
6. **Authentication** - može se dodati centralizovana autentifikacija

## Komponente

### API Gateway (`/api-gateway/`)
- **main.go** - glavna aplikacija
- **routes.go** - definisanje ruta za sve servise  
- **middleware.go** - CORS i logging middleware
- **Dockerfile** - Docker build konfiguracija

### Frontend (`/frontend/tour-app/`)
- **proxy.conf.json** - Angular dev proxy konfiguracija
- **nginx.conf** - produkcijska Nginx konfiguracija
- **environments/** - environment varijable za dev/prod
- **services/** - Angular servisi koji koriste API Gateway

Ova implementacija omogućava potpuno izolovan i centralizovan pristup komunikaciji između frontend-a i backend mikroservisa.

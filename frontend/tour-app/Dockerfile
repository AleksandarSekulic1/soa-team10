# Dockerfile za Angular frontend
FROM node:18-alpine AS builder

WORKDIR /app

# Kopiraj package.json i package-lock.json
COPY package*.json ./

# Instaliraj dependencies
RUN npm ci

# Kopiraj source kod
COPY . .

# Kopiraj Leaflet ikone u assets
RUN mkdir -p src/assets/images && \
    cp node_modules/leaflet/dist/images/marker-icon.png src/assets/images/ && \
    cp node_modules/leaflet/dist/images/marker-icon-2x.png src/assets/images/ && \
    cp node_modules/leaflet/dist/images/marker-shadow.png src/assets/images/

# Build aplikaciju za produkciju
RUN npm run build

# Production stage sa nginx
FROM nginx:alpine

# Kopiraj build fajlove
COPY --from=builder /app/dist/tour-app /usr/share/nginx/html

# Kopiraj nginx konfiguraciju
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

#!/bin/bash

# Test script za kreiranje korisnika i follow veza
API_BASE="http://localhost:8080"

echo "🚀 Kreiranje test korisnika za follower funkcionalnost..."

# Kreiraj test korisnike u stakeholders servisu
echo "📝 Kreiranje korisnika..."

# Korisnik 1: Ana (ID: 1)
curl -X POST "${API_BASE}/api/stakeholders/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ana_blogger",
    "email": "ana@example.com",
    "password": "password123",
    "role": "turista",
    "firstName": "Ana",
    "lastName": "Marković"
  }' && echo

# Korisnik 2: Milan (ID: 2)  
curl -X POST "${API_BASE}/api/stakeholders/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "milan_writer",
    "email": "milan@example.com", 
    "password": "password123",
    "role": "vodic",
    "firstName": "Milan",
    "lastName": "Petrović"
  }' && echo

# Korisnik 3: Marija (ID: 3)
curl -X POST "${API_BASE}/api/stakeholders/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "marija_travel",
    "email": "marija@example.com",
    "password": "password123", 
    "role": "turista",
    "firstName": "Marija",
    "lastName": "Nikolić"
  }' && echo

echo "✅ Korisnici kreirani!"

# Sačekaj malo
sleep 2

echo "🔗 Kreiranje korisnika u follower servisu..."

# Kreiraj korisnike u follower servisu
curl -X POST "${API_BASE}/api/followers/api/users" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "username": "ana_blogger",
    "email": "ana@example.com",
    "firstName": "Ana",
    "lastName": "Marković"
  }' && echo

curl -X POST "${API_BASE}/api/followers/api/users" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "username": "milan_writer", 
    "email": "milan@example.com",
    "firstName": "Milan",
    "lastName": "Petrović"
  }' && echo

curl -X POST "${API_BASE}/api/followers/api/users" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 3,
    "username": "marija_travel",
    "email": "marija@example.com",
    "firstName": "Marija", 
    "lastName": "Nikolić"
  }' && echo

echo "✅ Korisnici u follower servisu kreirani!"

# Sačekaj malo
sleep 2

echo "👥 Kreiranje follow veza..."

# Ana (1) prati Milana (2)
curl -X POST "${API_BASE}/api/followers/api/follow" \
  -H "Content-Type: application/json" \
  -d '{
    "followerId": 1,
    "followingId": 2
  }' && echo

# Ana (1) prati Mariju (3)  
curl -X POST "${API_BASE}/api/followers/api/follow" \
  -H "Content-Type: application/json" \
  -d '{
    "followerId": 1,
    "followingId": 3
  }' && echo

# Milan (2) prati Anu (1)
curl -X POST "${API_BASE}/api/followers/api/follow" \
  -H "Content-Type: application/json" \
  -d '{
    "followerId": 2,
    "followingId": 1
  }' && echo

# Marija (3) prati Anu (1)
curl -X POST "${API_BASE}/api/followers/api/follow" \
  -H "Content-Type: application/json" \
  -d '{
    "followerId": 3,
    "followingId": 1
  }' && echo

echo "✅ Follow veze kreirane!"

echo "🎉 Test podaci uspešno kreirani!"
echo ""
echo "📊 Test podaci:"
echo "👤 Ana (ID: 1) - prati: Milan, Mariju"
echo "👤 Milan (ID: 2) - prati: Anu"  
echo "👤 Marija (ID: 3) - prati: Anu"
echo ""
echo "💡 Sada se uloguj kao bilo koji korisnik da testiraš funkcionalnost!"
echo "   Username/Password kombinacije:"
echo "   • ana_blogger / password123"
echo "   • milan_writer / password123" 
echo "   • marija_travel / password123"
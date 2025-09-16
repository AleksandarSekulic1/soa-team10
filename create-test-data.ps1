# Test script za kreiranje korisnika i follow veza - PowerShell verzija
$API_BASE = "http://localhost:8080"

Write-Host "🚀 Kreiranje test korisnika za follower funkcionalnost..." -ForegroundColor Green

# Kreiraj test korisnike u stakeholders servisu
Write-Host "📝 Kreiranje korisnika..." -ForegroundColor Yellow

# Korisnik 1: Ana (ID treba da bude sledeći dostupan)
$ana = @{
    username = "ana_blogger"
    email = "ana@example.com"
    password = "password123"
    role = "turista"
    firstName = "Ana"
    lastName = "Marković"
} | ConvertTo-Json

try {
    $response1 = Invoke-RestMethod -Uri "$API_BASE/api/stakeholders/register" -Method POST -Body $ana -ContentType "application/json"
    Write-Host "✅ Ana kreirana" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Ana možda već postoji" -ForegroundColor Yellow
}

# Korisnik 2: Milan
$milan = @{
    username = "milan_writer"
    email = "milan@example.com"
    password = "password123"
    role = "vodic"
    firstName = "Milan"
    lastName = "Petrović"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri "$API_BASE/api/stakeholders/register" -Method POST -Body $milan -ContentType "application/json"
    Write-Host "✅ Milan kreiran" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Milan možda već postoji" -ForegroundColor Yellow
}

# Korisnik 3: Marija
$marija = @{
    username = "marija_travel"
    email = "marija@example.com"
    password = "password123"
    role = "turista"
    firstName = "Marija"
    lastName = "Nikolić"
} | ConvertTo-Json

try {
    $response3 = Invoke-RestMethod -Uri "$API_BASE/api/stakeholders/register" -Method POST -Body $marija -ContentType "application/json"
    Write-Host "✅ Marija kreirana" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Marija možda već postoji" -ForegroundColor Yellow
}

Start-Sleep -Seconds 2

Write-Host "🔗 Kreiranje korisnika u follower servisu..." -ForegroundColor Yellow

# Da dobijemo stvarne ID-jeve korisnika
try {
    $users = Invoke-RestMethod -Uri "$API_BASE/api/stakeholders" -Method GET
    
    $anaUser = $users | Where-Object { $_.username -eq "ana_blogger" } | Select-Object -First 1
    $milanUser = $users | Where-Object { $_.username -eq "milan_writer" } | Select-Object -First 1
    $marijaUser = $users | Where-Object { $_.username -eq "marija_travel" } | Select-Object -First 1
    
    if ($anaUser) {
        $anaFollower = @{
            id = [int]$anaUser.id
            username = $anaUser.username
            email = $anaUser.email
            firstName = $anaUser.firstName
            lastName = $anaUser.lastName
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/users" -Method POST -Body $anaFollower -ContentType "application/json"
            Write-Host "✅ Ana u follower servisu" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Ana već u follower servisu" -ForegroundColor Yellow
        }
    }
    
    if ($milanUser) {
        $milanFollower = @{
            id = [int]$milanUser.id
            username = $milanUser.username
            email = $milanUser.email
            firstName = $milanUser.firstName
            lastName = $milanUser.lastName
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/users" -Method POST -Body $milanFollower -ContentType "application/json"
            Write-Host "✅ Milan u follower servisu" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Milan već u follower servisu" -ForegroundColor Yellow
        }
    }
    
    if ($marijaUser) {
        $marijaFollower = @{
            id = [int]$marijaUser.id
            username = $marijaUser.username
            email = $marijaUser.email
            firstName = $marijaUser.firstName
            lastName = $marijaUser.lastName
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/users" -Method POST -Body $marijaFollower -ContentType "application/json"
            Write-Host "✅ Marija u follower servisu" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Marija već u follower servisu" -ForegroundColor Yellow
        }
    }
    
    Start-Sleep -Seconds 2
    
    Write-Host "👥 Kreiranje follow veza..." -ForegroundColor Yellow
    
    # Ana prati Milana
    if ($anaUser -and $milanUser) {
        $follow1 = @{
            followerId = [int]$anaUser.id
            followingId = [int]$milanUser.id
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/follow" -Method POST -Body $follow1 -ContentType "application/json"
            Write-Host "✅ Ana prati Milana" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Ana već prati Milana" -ForegroundColor Yellow
        }
    }
    
    # Ana prati Mariju
    if ($anaUser -and $marijaUser) {
        $follow2 = @{
            followerId = [int]$anaUser.id
            followingId = [int]$marijaUser.id
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/follow" -Method POST -Body $follow2 -ContentType "application/json"
            Write-Host "✅ Ana prati Mariju" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Ana već prati Mariju" -ForegroundColor Yellow
        }
    }
    
    # Milan prati Anu
    if ($milanUser -and $anaUser) {
        $follow3 = @{
            followerId = [int]$milanUser.id
            followingId = [int]$anaUser.id
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/follow" -Method POST -Body $follow3 -ContentType "application/json"
            Write-Host "✅ Milan prati Anu" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Milan već prati Anu" -ForegroundColor Yellow
        }
    }
    
    # Marija prati Anu
    if ($marijaUser -and $anaUser) {
        $follow4 = @{
            followerId = [int]$marijaUser.id
            followingId = [int]$anaUser.id
        } | ConvertTo-Json
        
        try {
            Invoke-RestMethod -Uri "$API_BASE/api/followers/api/follow" -Method POST -Body $follow4 -ContentType "application/json"
            Write-Host "✅ Marija prati Anu" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  Marija već prati Anu" -ForegroundColor Yellow
        }
    }
    
    Write-Host ""
    Write-Host "🎉 Test podaci uspešno kreirani!" -ForegroundColor Green
    Write-Host ""
    Write-Host "📊 Test podaci:" -ForegroundColor Cyan
    Write-Host "👤 Ana (ID: $($anaUser.id)) - prati: Milan, Mariju" -ForegroundColor White
    Write-Host "👤 Milan (ID: $($milanUser.id)) - prati: Anu" -ForegroundColor White
    Write-Host "👤 Marija (ID: $($marijaUser.id)) - prati: Anu" -ForegroundColor White
    Write-Host ""
    Write-Host "💡 Sada se uloguj kao bilo koji korisnik da testiraš funkcionalnost!" -ForegroundColor Yellow
    Write-Host "   Username/Password kombinacije:" -ForegroundColor Yellow
    Write-Host "   • ana_blogger / password123" -ForegroundColor White
    Write-Host "   • milan_writer / password123" -ForegroundColor White
    Write-Host "   • marija_travel / password123" -ForegroundColor White
    
} catch {
    Write-Host "❌ Greška pri dobijanju korisnika: $($_.Exception.Message)" -ForegroundColor Red
}
}

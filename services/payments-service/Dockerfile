# Faza 1: Build faza
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Kopiraj .csproj i restoruj zavisnosti
COPY *.csproj .
RUN dotnet restore

# Kopiraj ostatak koda i build-uj aplikaciju
COPY . .
RUN dotnet publish -c Release -o /app/publish

# Faza 2: Finalna faza
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app
COPY --from=build /app/publish .

# .NET 9 podrazumevano sluša na portu 8080
EXPOSE 8080

# ISPRAVLJENA LINIJA: Koristimo ispravno ime DLL fajla
ENTRYPOINT ["dotnet", "payments-service.dll"]
package service

import (
	"bytes"
	"encoding/json"
	"errors" // <-- 1. DODAT IMPORT
	"log"
	"net/http"
	"time"
	"tours-service/domain"
	"tours-service/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
// --- NOVO: Unesite vaš API ključ ovde ---
const orsApiKey = "eyJvcmciOiI1YjNjZTM1OTc4NTExMTAwMDFjZjYyNDgiLCJpZCI6IjQxMTE3OTdhMmE5YjQxM2Y5OWQzZWQ4MzM2MmNjY2VjIiwiaCI6Im11cm11cjY0In0="

// --- NOVO: Strukture za čitanje odgovora od ORS API-ja ---
type orsRequestBody struct {
	Coordinates [][]float64 `json:"coordinates"`
}
type orsSummary struct {
	Distance float64 `json:"distance"` // U metrima
}
type orsProperties struct {
	Summary orsSummary `json:"summary"`
}
type orsFeature struct {
	Properties orsProperties `json:"properties"`
}
type orsResponseBody struct {
	Features []orsFeature `json:"features"`
}

type TourService interface {
	Create(tour *domain.Tour) error
	GetById(tourId string) (*domain.Tour, error) // <-- DODATI
	GetByAuthorId(authorId string) ([]*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)                          // <-- DODATO
	AddReview(tourId string, review *domain.TourReview) error // <-- DODATO
	// NOVE METODE ZA KEY POINTS
	AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	DeleteKeyPoint(tourId, keyPointId string) error
    Update(tour *domain.Tour) (*domain.Tour, error) // <-- ISPRAVLJEN POVRATNI TIP
	// --- NOVE METODE ZA STANJA TURE ---
	Publish(tourId string) (*domain.Tour, error)
	Archive(tourId string) (*domain.Tour, error)
	Reactivate(tourId string) (*domain.Tour, error) // <-- DODATA NOVA METODA
	AddTransportInfo(tourId string, transportInfo []domain.TourTransport) (*domain.Tour, error)
	GetPublished() ([]*domain.Tour, error) // <-- DODATA NOVA METODA
	GetArchived() ([]*domain.Tour, error) // NOVO
}

type tourService struct {
	repo       repository.TourRepository
	httpClient *http.Client // <-- NOVO: HTTP Klijent za slanje zahteva
}

func NewTourService(repo repository.TourRepository) TourService {
	return &tourService{
		repo:       repo,
		httpClient: &http.Client{}, // <-- NOVO: Inicijalizacija klijenta
	}
}

// service/tour_service.go

func (s *tourService) Create(tour *domain.Tour) error {
	tour.Status = domain.TourStatusDraft // Koristimo konstantu
	tour.Price = 0.0
	tour.KeyPoints = []domain.TourKeyPoint{}
	tour.Reviews = []domain.TourReview{}
	tour.TransportInfo = []domain.TourTransport{}
	return s.repo.Create(tour)
}

func (s *tourService) GetById(tourId string) (*domain.Tour, error) {
	return s.repo.GetById(tourId)
}

func (s *tourService) GetByAuthorId(authorId string) ([]*domain.Tour, error) {
	return s.repo.GetByAuthorId(authorId)
}

func (s *tourService) GetAll() ([]*domain.Tour, error) {
	return s.repo.GetAll()
}

func (s *tourService) AddReview(tourId string, review *domain.TourReview) error {
	// Postavljamo serverske vrednosti
	review.ID = primitive.NewObjectID()
	review.CommentDate = time.Now()
	return s.repo.AddReview(tourId, review)
}

// --- PREPRAVLJENA AddKeyPoint METODA ---
func (s *tourService) AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	keyPoint.ID = primitive.NewObjectID()
	keyPoint.TourId = tourId
	
	// 1. Dodajemo ključnu tačku u bazu
	err := s.repo.AddKeyPoint(tourId, keyPoint)
	if err != nil {
		return err
	}

	// 2. Dobavljamo celu turu sa novom tačkom
	tour, err := s.repo.GetById(tourId)
	if err != nil {
		return err // Vraćamo grešku ako ne možemo da dobavimo turu
	}

	// 3. Ako ima 2 ili više tačaka, računamo distancu
	if len(tour.KeyPoints) >= 2 {
		log.Println("Calculating distance for tour:", tour.ID.Hex())
		
		// Pripremamo koordinate za ORS API ([longitude, latitude])
		var coords [][]float64
		for _, kp := range tour.KeyPoints {
			coords = append(coords, []float64{kp.Longitude, kp.Latitude})
		}

		// Pravimo telo zahteva
		reqBody := orsRequestBody{Coordinates: coords}
		reqBytes, _ := json.Marshal(reqBody)

		// Šaljemo POST zahtev na ORS API
		orsURL := "https://api.openrouteservice.org/v2/directions/driving-car/geojson"
		req, _ := http.NewRequest("POST", orsURL, bytes.NewBuffer(reqBytes))
		req.Header.Set("Authorization", orsApiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			log.Printf("Error calling ORS API: %v\n", err)
			return nil // Opciono: ne vraćamo grešku klijentu ako ORS ne radi
		}
		defer resp.Body.Close()

		// Čitamo i parsiramo odgovor
		var orsResp orsResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&orsResp); err != nil {
			log.Printf("Error decoding ORS response: %v\n", err)
			return nil
		}

		// Ažuriramo distancu na turi
		if len(orsResp.Features) > 0 {
			distanceInMeters := orsResp.Features[0].Properties.Summary.Distance
			tour.Distance = distanceInMeters / 1000 // Pretvaramo u km
			log.Printf("New distance for tour %s is %f km\n", tour.ID.Hex(), tour.Distance)
			
			// Sačuvamo izmenjenu turu u bazi
			return s.repo.Update(tour)
		}
	}

	return nil
}

func (s *tourService) UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	// Logika za validaciju ili dodatne provere bi išla ovde
	return s.repo.UpdateKeyPoint(tourId, keyPoint)
}

func (s *tourService) DeleteKeyPoint(tourId, keyPointId string) error {
	return s.repo.DeleteKeyPoint(tourId, keyPointId)
}

// --- DODAJEMO NOVU METODU ---
func (s *tourService) Update(tour *domain.Tour) (*domain.Tour, error) {
	err := s.repo.Update(tour)
	if err != nil {
		return nil, err
	}
	// Vraćamo ažuriranu turu ponovnim dobavljanjem iz baze
	return s.repo.GetById(tour.ID.Hex())
}
func (s *tourService) AddTransportInfo(tourId string, transportInfo []domain.TourTransport) (*domain.Tour, error) {
	tour, err := s.repo.GetById(tourId)
	if err != nil {
		return nil, err
	}
	tour.TransportInfo = transportInfo
	return s.Update(tour)
}

func (s *tourService) Publish(tourId string) (*domain.Tour, error) {
	tour, err := s.repo.GetById(tourId)
	if err != nil {
		return nil, err
	}

	// Provera 1: Tura mora biti u 'draft' statusu
	if tour.Status != domain.TourStatusDraft {
		return nil, errors.New("tour can only be published if it is in draft status")
	}

	// Provera 2: Tura mora imati bar dve ključne tačke
	if len(tour.KeyPoints) < 2 {
		return nil, errors.New("tour must have at least two key points to be published")
	}

	// Provera 3: Tura mora imati bar jednu informaciju o transportu
	if len(tour.TransportInfo) < 1 {
		return nil, errors.New("tour must have at least one transport time defined")
	}

	// Ako su svi uslovi ispunjeni, menjamo status
	now := time.Now()
	tour.Status = domain.TourStatusPublished
	tour.PublishedAt = &now

	return s.Update(tour)
}

func (s *tourService) Archive(tourId string) (*domain.Tour, error) {
	tour, err := s.repo.GetById(tourId)
	if err != nil {
		return nil, err
	}

	// Tura se može arhivirati samo ako je objavljena
	if tour.Status != domain.TourStatusPublished {
		return nil, errors.New("tour can only be archived if it is in published status")
	}

	now := time.Now()
	tour.Status = domain.TourStatusArchived
	tour.ArchivedAt = &now

	return s.Update(tour)
}

func (s *tourService) Reactivate(tourId string) (*domain.Tour, error) {
	tour, err := s.repo.GetById(tourId)
	if err != nil {
		return nil, err
	}

	// Tura se može reaktivirati samo ako je arhivirana
	if tour.Status != domain.TourStatusArchived {
		return nil, errors.New("tour can only be reactivated if it is in archived status")
	}

	tour.Status = domain.TourStatusPublished
	tour.ArchivedAt = nil // Poništavamo vreme arhiviranja

	return s.Update(tour)
}

func (s *tourService) GetPublished() ([]*domain.Tour, error) {
	// 1. Dobavi samo objavljene ture iz repozitorijuma
	publishedTours, err := s.repo.GetPublished()
	if err != nil {
		return nil, err
	}

	// 2. Za svaku turu, ostavi samo prvu ključnu tačku
	for _, tour := range publishedTours {
		if len(tour.KeyPoints) > 0 {
			tour.KeyPoints = tour.KeyPoints[:1] // Skraćujemo niz na samo prvi element
		}
	}

	return publishedTours, nil
}

// NOVO: Servisna metoda za dobavljanje arhiviranih tura
func (s *tourService) GetArchived() ([]*domain.Tour, error) {
	// 1. Dobavi samo arhivirane ture iz repozitorijuma
	archivedTours, err := s.repo.GetArchived()
	if err != nil {
		return nil, err
	}

	// 2. Za svaku turu, ostavi samo prvu ključnu tačku (kao i za objavljene)
	for _, tour := range archivedTours {
		if len(tour.KeyPoints) > 0 {
			tour.KeyPoints = tour.KeyPoints[:1] // Skraćujemo niz na samo prvi element
		}
	}

	return archivedTours, nil
}

package api

import (
	"log"
	"net/http"
	"tours-service/domain"
	"tours-service/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourHandler struct {
	service service.TourService
}

func NewTourHandler(service service.TourService) *TourHandler {
	return &TourHandler{service: service}
}

// @Summary Kreiranje nove ture
// @Description Kreira novu turu sa statusom 'draft' i cenom 0. Zahteva autentikaciju.
// @Accept  json
// @Produce  json
// @Param   tour body domain.Tour true "Podaci o turi (bez statusa, cene i authorId)"
// @Security ApiKeyAuth
// @Success 201 {object} domain.Tour "Uspešno kreirana tura"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Router /tours [post]
func (h *TourHandler) Create(c *gin.Context) {
	var tour domain.Tour
	if err := c.ShouldBindJSON(&tour); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}
	authorUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}
	tour.AuthorId = authorUsername.(string)
	err := h.service.Create(&tour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom kreiranja ture"})
		return
	}
	c.JSON(http.StatusCreated, tour)
}

// @Summary Prikaz tura kreiranih od strane ulogovanog autora
// @Description Vraća listu svih tura koje je kreirao autor čiji se token koristi.
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} domain.Tour "Lista tura"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Router /tours/my-tours [get]
func (h *TourHandler) GetByAuthor(c *gin.Context) {
	authorUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}
	tours, err := h.service.GetByAuthorId(authorUsername.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja tura"})
		return
	}
	c.JSON(http.StatusOK, tours)
}

// @Summary Prikaz svih tura
// @Description Vraća listu svih tura dostupnih u sistemu.
// @Produce  json
// @Success 200 {array} domain.Tour "Lista svih tura"
// @Router /tours [get]
func (h *TourHandler) GetAll(c *gin.Context) {
	tours, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja tura"})
		return
	}
	c.JSON(http.StatusOK, tours)
}

// @Summary Dodavanje recenzije na turu
// @Description Dodaje novu recenziju na turu sa datim ID-jem. Zahteva autentikaciju.
// @Accept  json
// @Produce  json
// @Param   id   path      string  true  "ID Ture"
// @Param   review body domain.TourReview true "Podaci o recenziji"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Router /tours/{id}/reviews [post]
func (h *TourHandler) AddReview(c *gin.Context) {
	tourId := c.Param("id")
	var review domain.TourReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}
	touristUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}
	review.TouristId = touristUsername.(string)
	err := h.service.AddReview(tourId, &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom dodavanja recenzije"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recenzija je uspešno dodata."})
}

// @Summary Dodavanje ključne tačke na turu
// @Description Dodaje novu ključnu tačku na turu sa datim ID-jem. Samo autor ture može dodati tačku.
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Param   keyPoint body domain.TourKeyPoint true "Podaci o ključnoj tački"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TourKeyPoint "Uspešno dodata ključna tačka"
// @Router /tours/{id}/keypoints [post]
func (h *TourHandler) AddKeyPoint(c *gin.Context) {
	tourId := c.Param("id")
	var keyPoint domain.TourKeyPoint
	if err := c.ShouldBindJSON(&keyPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}
	
	// --- POCETAK PROVERE AUTORIZACIJE ---
	authorUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}

	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tura nije pronađena"})
		return
	}

	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nemate dozvolu da menjate ovu turu"})
		return
	}
	// --- KRAJ PROVERE AUTORIZACIJE ---

	if err := h.service.AddKeyPoint(tourId, &keyPoint); err != nil {
		log.Printf("!!! SERVER ERROR - AddKeyPoint: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom dodavanja ključne tačke"})
		return
	}
	c.JSON(http.StatusOK, keyPoint)
}

// @Summary Ažuriranje ključne tačke
// @Description Ažurira postojeću ključnu tačku na turi.
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Param   keypointId path string true "ID Ključne tačke"
// @Param   keyPoint body domain.TourKeyPoint true "Novi podaci o ključnoj tački"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TourKeyPoint "Uspešno ažurirana ključna tačka"
// @Router /tours/{id}/keypoints/{keypointId} [put]
func (h *TourHandler) UpdateKeyPoint(c *gin.Context) {
	tourId := c.Param("id")
	keyPointId := c.Param("keypointId")

	var keyPoint domain.TourKeyPoint
	if err := c.ShouldBindJSON(&keyPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}

	// --- POCETAK PROVERE AUTORIZACIJE ---
	authorUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}

	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tura nije pronađena"})
		return
	}

	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nemate dozvolu da menjate ovu turu"})
		return
	}
	// --- KRAJ PROVERE AUTORIZACIJE ---

	keyPointObjID, err := primitive.ObjectIDFromHex(keyPointId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan ID format za ključnu tačku"})
		return
	}
	keyPoint.ID = keyPointObjID
	keyPoint.TourId = tourId

	if err := h.service.UpdateKeyPoint(tourId, &keyPoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom ažuriranja ključne tačke"})
		return
	}
	c.JSON(http.StatusOK, keyPoint)
}

// @Summary Brisanje ključne tačke
// @Description Briše ključnu tačku sa ture.
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Param   keypointId path string true "ID Ključne tačke"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Router /tours/{id}/keypoints/{keypointId} [delete]
func (h *TourHandler) DeleteKeyPoint(c *gin.Context) {
	tourId := c.Param("id")
	keyPointId := c.Param("keypointId")

	// --- POCETAK PROVERE AUTORIZACIJE ---
	authorUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autorizovan"})
		return
	}

	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tura nije pronađena"})
		return
	}

	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nemate dozvolu da menjate ovu turu"})
		return
	}
	// --- KRAJ PROVERE AUTORIZACIJE ---

	if err := h.service.DeleteKeyPoint(tourId, keyPointId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom brisanja ključne tačke"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ključna tačka je uspešno obrisana."})
}

// @Summary Prikaz jedne ture po ID-ju
// @Description Vraća detalje specifične ture na osnovu njenog ID-ja.
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Success 200 {object} domain.Tour "Detalji ture"
// @Failure 404 {object} map[string]string "Greška: Tura nije pronađena"
// @Router /tours/{id} [get]
func (h *TourHandler) GetById(c *gin.Context) {
	tourId := c.Param("id")

	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tura nije pronađena"})
		return
	}

	c.JSON(http.StatusOK, tour)
}

// @Summary Dodavanje informacija o transportu
// @Description Dodaje ili ažurira informacije o vremenu putovanja za turu.
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Param   transportInfo  body  []domain.TourTransport  true  "Lista informacija o transportu"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Tour "Uspešno ažurirana tura"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 403 {object} map[string]string "Greška: Nemate dozvolu"
// @Router /tours/{id}/transport-info [post]
func (h *TourHandler) AddTransportInfo(c *gin.Context) {
	tourId := c.Param("id")
	var transportInfo []domain.TourTransport
	if err := c.ShouldBindJSON(&transportInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Provera autorizacije: Samo autor može menjati turu
	authorUsername, _ := c.Get("username")
	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}
	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this tour"})
		return
	}

	updatedTour, err := h.service.AddTransportInfo(tourId, transportInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTour)
}

// @Summary Objavljivanje ture
// @Description Menja status ture u 'published' ako su ispunjeni svi uslovi.
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Tour "Uspešno objavljena tura"
// @Failure 400 {object} map[string]string "Greška: Uslovi za objavljivanje nisu ispunjeni"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 403 {object} map[string]string "Greška: Nemate dozvolu"
// @Router /tours/{id}/publish [post]
func (h *TourHandler) Publish(c *gin.Context) {
	tourId := c.Param("id")

	// Provera autorizacije
	authorUsername, _ := c.Get("username")
	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}
	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this tour"})
		return
	}

	publishedTour, err := h.service.Publish(tourId)
	if err != nil {
		// Vraćamo Bad Request jer uslovi nisu ispunjeni
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publishedTour)
}

// @Summary Arhiviranje ture
// @Description Menja status ture u 'archived'.
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Tour "Uspešno arhivirana tura"
// @Failure 400 {object} map[string]string "Greška: Tura nije objavljena"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 403 {object} map[string]string "Greška: Nemate dozvolu"
// @Router /tours/{id}/archive [post]
func (h *TourHandler) Archive(c *gin.Context) {
	tourId := c.Param("id")

	// Provera autorizacije
	authorUsername, _ := c.Get("username")
	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}
	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this tour"})
		return
	}

	archivedTour, err := h.service.Archive(tourId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, archivedTour)
}

// @Summary Prikaz svih objavljenih tura (za turiste)
// @Description Vraća listu tura koje imaju status 'published'. Za svaku turu prikazuje samo prvu ključnu tačku.
// @Produce  json
// @Success 200 {array} domain.Tour "Lista objavljenih tura"
// @Router /tours/published [get]
func (h *TourHandler) GetPublished(c *gin.Context) {
	tours, err := h.service.GetPublished()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja objavljenih tura"})
		return
	}
	c.JSON(http.StatusOK, tours)
}

// @Summary Ponovno aktiviranje ture
// @Description Menja status arhivirane ture nazad u 'published'.
// @Produce  json
// @Param   id   path   string  true  "ID Ture"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Tour "Uspešno reaktivirana tura"
// @Failure 400 {object} map[string]string "Greška: Tura nije arhivirana"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 403 {object} map[string]string "Greška: Nemate dozvolu"
// @Router /tours/{id}/reactivate [post]
func (h *TourHandler) Reactivate(c *gin.Context) {
	tourId := c.Param("id")

	// Provera autorizacije: Samo autor može menjati turu
	authorUsername, _ := c.Get("username")
	tour, err := h.service.GetById(tourId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}
	if tour.AuthorId != authorUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this tour"})
		return
	}

	reactivatedTour, err := h.service.Reactivate(tourId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reactivatedTour)
}

// @Summary Prikaz svih arhiviranih tura
// @Description Vraća listu tura koje imaju status 'archived'. Za svaku turu prikazuje samo prvu ključnu tačku.
// @Produce 	json
// @Success 200 {array} domain.Tour "Lista arhiviranih tura"
// @Router /tours/archived [get]
func (h *TourHandler) GetArchived(c *gin.Context) {
	tours, err := h.service.GetArchived()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja arhiviranih tura"})
		return
	}
	c.JSON(http.StatusOK, tours)
}
package api

import (
	"net/http"
	"tours-service/domain"
	"tours-service/service"

	"github.com/gin-gonic/gin"
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
	// Preuzimamo korisničko ime iz tokena koje je postavio middleware
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

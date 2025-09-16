// api/tourist_position_handler.go

package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encounters-service/domain"
	"encounters-service/service"
)

type TouristPositionHandler struct {
	service service.TouristPositionService
}

func NewTouristPositionHandler(service service.TouristPositionService) *TouristPositionHandler {
	return &TouristPositionHandler{service: service}
}

type PositionUpdateDTO struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// @Summary Dobavljanje pozicije turiste
// @Description Vraća poslednju sačuvanu geografsku poziciju za ulogovanog korisnika.
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} domain.TouristPosition "Uspešno dobavljena pozicija"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 404 {object} map[string]string "Greška: Pozicija nije pronađena"
// @Router /tourist-position [get]
func (h *TouristPositionHandler) GetByUserId(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	position, err := h.service.GetByUserId(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found for user"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// @Summary Ažuriranje pozicije turiste
// @Description Čuva ili ažurira geografsku poziciju za ulogovanog korisnika. Koristi se za simulator pozicije.
// @Accept  json
// @Produce  json
// @Param   position  body  PositionUpdateDTO  true  "Podaci o novoj poziciji (Latitude i Longitude)"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TouristPosition "Uspešno ažurirana pozicija"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 500 {object} map[string]string "Greška: Problem na serveru"
// @Router /tourist-position [post]
func (h *TouristPositionHandler) Update(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	var dto PositionUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	position := &domain.TouristPosition{
		UserId:    username.(string),
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
	}

	updatedPosition, err := h.service.Update(position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	c.JSON(http.StatusOK, updatedPosition)
}
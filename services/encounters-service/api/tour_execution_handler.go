// encounters-service/api/tour_execution_handler.go
package api

import (
	"encounters-service/domain"
	"encounters-service/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive" // <-- DODAT IMPORT
	"net/http"
	// strconv je uklonjen
)

type TourExecutionHandler struct {
	service service.TourExecutionService
}

func NewTourExecutionHandler(service service.TourExecutionService) *TourExecutionHandler {
	return &TourExecutionHandler{service: service}
}

type CheckPositionDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// @Summary Pokretanje ture
// @Description Kreira novu sesiju izvođenja ture za ulogovanog korisnika.
// @Accept  json
// @Produce  json
// @Param   tourId  path  string  true  "ID Ture koja se pokreće"
// @Security ApiKeyAuth
// @Success 201 {object} domain.TourExecution "Uspešno kreirana sesija"
// @Failure 400 {object} map[string]string "Greška: Korisnik već ima aktivnu turu ili neispravan ID"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Router /tour-executions/start/{tourId} [post]
func (h *TourExecutionHandler) StartTour(c *gin.Context) {
	tourId := c.Param("tourId")

	// Proveravamo da li je ID validan ObjectID, ako želimo dodatnu sigurnost
	if _, err := primitive.ObjectIDFromHex(tourId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID format"})
		return
	}

	userId, _ := c.Get("userId")

	execution := &domain.TourExecution{
		TourId: tourId, // Sada se tipovi poklapaju (string u string)
		UserId: userId.(string),
	}

	newExecution, err := h.service.StartTour(execution)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newExecution)
}

// @Summary Provera pozicije
// @Description Proverava trenutnu poziciju turiste i ažurira napredak na aktivnoj turi.
// @Accept  json
// @Produce  json
// @Param   position  body  CheckPositionDTO  true  "Trenutne koordinate turiste"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TourExecution "Ažurirana sesija sa napretkom"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Failure 404 {object} map[string]string "Greška: Korisnik nema aktivnu turu"
// @Router /tour-executions/check-position [post]
func (h *TourExecutionHandler) CheckPosition(c *gin.Context) {
	userId, _ := c.Get("userId")

	var dto CheckPositionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedExecution, err := h.service.CheckPosition(userId.(string), dto.Latitude, dto.Longitude)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedExecution)
}

// @Summary Završetak ture
// @Description Označava aktivnu turu kao 'Completed'.
// @Accept  json
// @Produce  json
// @Param   executionId  path  string  true  "ID sesije koja se završava"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TourExecution "Završena sesija"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Router /tour-executions/{executionId}/complete [post]
func (h *TourExecutionHandler) CompleteTour(c *gin.Context) {
	executionId := c.Param("executionId")
	// TODO: Dodati proveru da li je korisnik vlasnik ove sesije
	updatedExecution, err := h.service.CompleteTour(executionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedExecution)
}

// @Summary Napuštanje ture
// @Description Označava aktivnu turu kao 'Abandoned'.
// @Accept  json
// @Produce  json
// @Param   executionId  path  string  true  "ID sesije koja se napušta"
// @Security ApiKeyAuth
// @Success 200 {object} domain.TourExecution "Napuštena sesija"
// @Failure 401 {object} map[string]string "Greška: Korisnik nije autorizovan"
// @Router /tour-executions/{executionId}/abandon [post]
func (h *TourExecutionHandler) AbandonTour(c *gin.Context) {
	executionId := c.Param("executionId")
	// TODO: Dodati proveru da li je korisnik vlasnik ove sesije
	updatedExecution, err := h.service.AbandonTour(executionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedExecution)
}
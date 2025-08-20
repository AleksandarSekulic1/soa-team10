// api/user_handler.go

package api

import (
	"net/http"
	"stakeholders-service/domain"
	"stakeholders-service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// @Summary Registracija novog korisnika
// @Description Kreira novog korisnika sa ulogom 'turista' ili 'vodič'.
// @Accept  json
// @Produce  json
// @Param   user body domain.User true "Podaci za registraciju korisnika"
// @Success 201 {object} domain.User "Uspešno kreiran korisnik"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}

	err := h.service.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom registracije korisnika"})
		return
	}

	// U pravoj aplikaciji ne vraćamo lozinku.
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

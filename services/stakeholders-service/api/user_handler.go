package api

import (
	"net/http"
	"stakeholders-service/domain"
	"stakeholders-service/service"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// @Summary Prikaz svih korisničkih naloga
// @Description Vraća listu svih korisničkih naloga bez lozinki.
// @Produce json
// @Success 200 {array} domain.User "Lista svih korisnika"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja korisnika"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Prijava korisnika
// @Description Autentifikuje korisnika i vraća JWT token.
// @Accept  json
// @Produce  json
// @Param   credentials body LoginRequest true "Korisničko ime i lozinka"
// @Success 200 {object} map[string]string "Uspešna prijava, vraća token"
// @Failure 401 {object} map[string]string "Greška: Neispravni kredencijali"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}
	token, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravno korisničko ime ili lozinka"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Interna greška servera"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

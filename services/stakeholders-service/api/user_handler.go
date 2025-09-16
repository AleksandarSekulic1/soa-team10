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

// ... Register, GetAll, Login ostaju isti sa svojim anotacijama ...
// @Summary Registracija novog korisnika
// @Description Kreira novog korisnika sa ulogom 'turista' ili 'vodič'.
// @Accept  json
// @Produce  json
// @Param   user body domain.User true "Podaci za registraciju korisnika"
// @Success 201 {object} domain.User "Uspešno kreiran korisnik"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    // ... postojeća Register metoda ostaje ista ...
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

// @Summary Prikaz svih korisničkih naloga (samo za admine)
// @Security ApiKeyAuth
// @Description Vraća listu svih korisničkih naloga bez lozinki.
// @Produce json
// @Success 200 {array} domain.User "Lista svih korisnika"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 403 {object} map[string]string "Pristup zabranjen (potrebna admin uloga)"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders [get]
func (h *UserHandler) GetAll(c *gin.Context) {
    // ... postojeća GetAll metoda ostaje ista ...
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom preuzimanja korisnika"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Prijava korisnika
// @Description Autentifikuje korisnika i vraća JWT token.
// @Accept  json
// @Produce  json
// @Param   credentials body LoginRequest true "Korisničko ime i lozinka"
// @Success 200 {object} map[string]string "Uspešna prijava, vraća token"
// @Failure 401 {object} map[string]string "Greška: Neispravni kredencijali"
// @Failure 403 {object} map[string]string "Korisnik je blokiran"
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
		// Provera za specifične greške iz servisa
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravno korisničko ime ili lozinka"})
			return
		}
		if err.Error() == "user is blocked" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Korisnik je blokiran"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Interna greška servera"})
		return
	}
	// Vraćamo samo token, kao što frontend očekuje
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// NOVA HANDLER METODA za dobavljanje profila
// @Summary Dobavljanje profila ulogovanog korisnika
// @Security ApiKeyAuth
// @Description Vraća sve podatke o profilu korisnika koji je trenutno ulogovan.
// @Produce json
// @Success 200 {object} domain.User "Podaci o profilu"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 404 {object} map[string]string "Korisnik nije pronađen"
// @Router /stakeholders/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Uzimamo username iz context-a koji je postavio middleware
	username, _ := c.Get("username")
	user, err := h.service.GetProfile(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Korisnik nije pronađen"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// NOVA HANDLER METODA za ažuriranje profila
// @Summary Ažuriranje profila ulogovanog korisnika
// @Security ApiKeyAuth
// @Description Ažurira podatke o profilu (ime, prezime, bio, moto...) za ulogovanog korisnika.
// @Accept json
// @Produce json
// @Param user body domain.User true "Podaci za ažuriranje profila"
// @Success 200 {object} domain.User "Ažurirani podaci o profilu"
// @Failure 400 {object} map[string]string "Neispravan format zahteva"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	username, _ := c.Get("username")
	var userUpdates domain.User
	if err := c.ShouldBindJSON(&userUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format zahteva"})
		return
	}
	// Osiguravamo da korisnik može menjati samo svoj profil
	userUpdates.Username = username.(string)

	updatedUser, err := h.service.UpdateProfile(&userUpdates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom ažuriranja profila"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// NOVA HANDLER METODA za blokiranje korisnika
// @Summary Blokiranje korisnika (samo za admine)
// @Security ApiKeyAuth
// @Description Postavlja status korisnika na 'blokiran'.
// @Produce json
// @Param username path string true "Korisničko ime korisnika koji se blokira"
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 403 {object} map[string]string "Pristup zabranjen (potrebna admin uloga)"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/{username}/block [put]
func (h *UserHandler) BlockUser(c *gin.Context) {
	username := c.Param("username")
	err := h.service.SetBlockedStatus(username, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom blokiranja korisnika"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Korisnik " + username + " je uspešno blokiran."})
}

// NOVA HANDLER METODA za odblokiranje korisnika
// @Summary Odblokiranje korisnika (samo za admine)
// @Security ApiKeyAuth
// @Description Postavlja status korisnika na 'aktivan'.
// @Produce json
// @Param username path string true "Korisničko ime korisnika koji se odblokira"
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 403 {object} map[string]string "Pristup zabranjen (potrebna admin uloga)"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /stakeholders/{username}/unblock [put]
func (h *UserHandler) UnblockUser(c *gin.Context) {
	username := c.Param("username")
	err := h.service.SetBlockedStatus(username, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom odblokiranja korisnika"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Korisnik " + username + " je uspešno odblokiran."})
}
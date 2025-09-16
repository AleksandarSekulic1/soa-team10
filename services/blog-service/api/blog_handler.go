// u services/blog-service/api/blog_handler.go

package api

import (
	"errors" // <-- DODAT IMPORT
	"net/http"

	"blog-service/domain"
	"blog-service/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struktura handlera
type BlogHandler struct {
	service service.BlogService
}

// Konstruktor za kreiranje novog handlera
func NewBlogHandler(service service.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

// @Summary Kreiranje novog bloga
// @Security ApiKeyAuth
// @Description Kreira novi blog post sa naslovom i sadržajem. Zahteva autorizaciju.
// @Accept  json
// @Produce  json
// @Param   blog body domain.Blog true "Podaci za kreiranje bloga (potrebni samo title i content)"
// @Success 201 {object} domain.Blog "Uspešno kreiran blog"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /blogs [post]
func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var blog domain.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// UZIMAMO PODATKE O AUTORU IZ TOKENA (KOJE JE POSTAVIO MIDDLEWARE)
	authorId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
		return
	}
	blog.AuthorID = authorId.(string) // Postavljamo pravog autora

	err := h.service.Create(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

// @Summary Prikaz svih objava na blogu
// @Security ApiKeyAuth
// @Description Vraća listu svih blog objava. Zahteva autorizaciju.
// @Produce  json
// @Success 200 {array} domain.Blog "Lista svih blogova"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /blogs [get]
func (h *BlogHandler) GetAllBlogs(c *gin.Context) {
	blogs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs"})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

// @Summary Dodavanje komentara na blog
// @Security ApiKeyAuth
// @Description Dodaje novi komentar na postojeći blog. Zahteva autorizaciju.
// @Accept  json
// @Produce  json
// @Param   id   path      string  true  "Blog ID"
// @Param   comment body domain.Comment true "Tekst komentara"
// @Success 200 {object} domain.Blog "Blog sa novim komentarom"
// @Failure 400 {object} map[string]string "Greška: Neispravan format zahteva ili ID-a"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /blogs/{id}/comments [post]
func (h *BlogHandler) AddComment(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var comment domain.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authorId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
		return
	}
	comment.AuthorID = authorId.(string)

	err = h.service.AddComment(blogID, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

// @Summary Lajkovanje/Uklanjanje lajka sa bloga
// @Security ApiKeyAuth
// @Description Dodaje ili uklanja lajk ulogovanog korisnika sa objave.
// @Produce  json
// @Param   id   path      string  true  "Blog ID"
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Failure 400 {object} map[string]string "Greška: Neispravan ID"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 500 {object} map[string]string "Interna greška servera"
// @Router /blogs/{id}/likes [post]
func (h *BlogHandler) ToggleLike(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	userID, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
		return
	}

	err = h.service.ToggleLike(blogID, userID.(string))
	if err != nil {
		// Ako je greška "mongo: no documents in result", vraćamo 404
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog with that ID not found"})
			return
		}
		// Za sve ostale greške, vraćamo 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like toggled successfully"})
}

// @Summary Prikaz jednog bloga
// @Security ApiKeyAuth
// @Description Vraća sve detalje o jednom blog postu na osnovu njegovog ID-a. Zahteva autorizaciju.
// @Produce  json
// @Param   id   path      string  true  "Blog ID"
// @Success 200 {object} domain.Blog "Detalji bloga"
// @Failure 400 {object} map[string]string "Greška: Neispravan ID"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 404 {object} map[string]string "Blog nije pronađen"
// @Router /blogs/{id} [get]
func (h *BlogHandler) GetBlogById(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	blog, err := h.service.GetById(blogID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// @Summary Ažuriranje bloga
// @Security ApiKeyAuth
// @Description Ažurira naslov ili sadržaj postojećeg bloga. Dozvoljeno samo autoru.
// @Accept  json
// @Produce  json
// @Param   id   path      string  true  "Blog ID"
// @Param   blog body domain.Blog true "Novi podaci za blog (title, content, images)"
// @Success 200 {object} domain.Blog "Ažurirani blog"
// @Failure 400 {object} map[string]string "Neispravan format zahteva ili ID-a"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 403 {object} map[string]string "Pristup zabranjen (niste autor)"
// @Failure 404 {object} map[string]string "Blog nije pronađen"
// @Router /blogs/{id} [put]
func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blogUpdates domain.Blog
	if err := c.ShouldBindJSON(&blogUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, _ := c.Get("username")

	updatedBlog, err := h.service.Update(blogID, &blogUpdates, userID.(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
			return
		}
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this blog"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, updatedBlog)
}

// @Summary Izmena komentara
// @Security ApiKeyAuth
// @Description Menja tekst postojećeg komentara. Dozvoljeno samo autoru komentara.
// @Accept  json
// @Produce  json
// @Param   id   path      string  true  "Blog ID"
// @Param   commentId   path      string  true  "Comment ID"
// @Param   comment body domain.Comment true "Novi tekst komentara (potrebno samo 'text' polje)"
// @Success 200 {object} map[string]string "Poruka o uspehu"
// @Failure 400 {object} map[string]string "Greška: Neispravan format ID-a ili zahteva"
// @Failure 401 {object} map[string]string "Neautorizovan pristup"
// @Failure 403 {object} map[string]string "Pristup zabranjen (niste autor)"
// @Failure 404 {object} map[string]string "Blog ili komentar nije pronađen"
// @Router /blogs/{id}/comments/{commentId} [put]
func (h *BlogHandler) UpdateComment(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var commentUpdate domain.Comment
	if err := c.ShouldBindJSON(&commentUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, _ := c.Get("username")

	err = h.service.UpdateComment(blogID, commentID, &commentUpdate, userID.(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
			return
		}
		if err.Error() == "comment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this comment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}
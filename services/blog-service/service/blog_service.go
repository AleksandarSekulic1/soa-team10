package service

import (
	"errors" // <-- DODAT IMPORT
	"time"

	"blog-service/domain"
	"blog-service/repository"

	"github.com/golang-jwt/jwt/v5" // <-- DODAJEMO IMPORT
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definišemo tajni ključ, mora biti isti kao u stakeholders-service
var JwtKey = []byte("super_secret_key")

// Definišemo Claims strukturu, mora biti ista kao u stakeholders-service
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}


type BlogService interface {
	Create(blog *domain.Blog) error
	GetAll() ([]*domain.Blog, error)
	GetBlogsByAuthors(authorIDs []string) ([]*domain.Blog, error) // <-- NOVA METODA
	AddComment(blogID primitive.ObjectID, comment *domain.Comment) error // <-- NOVA METODA
	ToggleLike(blogID primitive.ObjectID, userID string) error // <-- NOVA METODA
	GetById(id primitive.ObjectID) (*domain.Blog, error)       // <-- NOVA METODA
	Update(id primitive.ObjectID, blog *domain.Blog, userID string) (*domain.Blog, error) // <-- NOVA METODA
	UpdateComment(blogID, commentID primitive.ObjectID, updatedComment *domain.Comment, userID string) error // <-- NOVA METODA
	RemoveLikesFromAuthorBlogs(userID, authorID string) error // <-- NOVA METODA ZA SAGA
}

type blogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) Create(blog *domain.Blog) error {
	blog.ID = primitive.NewObjectID()
	blog.CreatedAt = time.Now()
	blog.Comments = []domain.Comment{}
	blog.Likes = []domain.Like{}
	
	return s.repo.Create(blog)
}

func (s *blogService) GetAll() ([]*domain.Blog, error) {
	return s.repo.GetAll()
}

// IMPLEMENTACIJA NOVE METODE ZA BLOGOVE OD PRAĆENIH KORISNIKA
func (s *blogService) GetBlogsByAuthors(authorIDs []string) ([]*domain.Blog, error) {
	if len(authorIDs) == 0 {
		// Ako nema praćenih korisnika, vrati praznu listu
		return []*domain.Blog{}, nil
	}
	return s.repo.GetBlogsByAuthors(authorIDs)
}

// IMPLEMENTACIJA NOVE METODE
func (s *blogService) AddComment(blogID primitive.ObjectID, comment *domain.Comment) error {
	// KLJUČNA ISPRAVKA: Generišemo novi jedinstveni ID za svaki komentar
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.LastUpdatedAt = time.Now()
	return s.repo.AddComment(blogID, comment)
}


// IMPLEMENTACIJA NOVE METODE
func (s *blogService) ToggleLike(blogID primitive.ObjectID, userID string) error {
	// 1. Prvo dobavljamo blog da vidimo da li lajk već postoji
	blog, err := s.repo.GetById(blogID)
	if err != nil {
		return err
	}

	// 2. Proveravamo da li userID postoji u listi lajkova
	for _, like := range blog.Likes {
		if like.UserID == userID {
			// Ako postoji, uklanjamo lajk i izlazimo iz funkcije
			return s.repo.RemoveLike(blogID, userID)
		}
	}

	// 3. Ako lajk nije pronađen, dodajemo novi
	newLike := &domain.Like{
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return s.repo.AddLike(blogID, newLike)
}

// IMPLEMENTACIJA NOVE METODE
func (s *blogService) GetById(id primitive.ObjectID) (*domain.Blog, error) {
	return s.repo.GetById(id)
}

// IMPLEMENTACIJA NOVE METODE
func (s *blogService) Update(id primitive.ObjectID, updatedBlog *domain.Blog, userID string) (*domain.Blog, error) {
	// 1. Proveravamo da li blog uopšte postoji
	existingBlog, err := s.repo.GetById(id)
	if err != nil {
		return nil, err // Vraća "not found" grešku ako ne postoji
	}

	// 2. Ključna provera: da li je korisnik koji menja blog i njegov autor
	if existingBlog.AuthorID != userID {
		return nil, errors.New("forbidden") // Vraćamo grešku ako nije autor
	}

	// 3. Ako jeste, ažuriramo podatke
	existingBlog.Title = updatedBlog.Title
	existingBlog.Content = updatedBlog.Content
	existingBlog.Images = updatedBlog.Images

	err = s.repo.Update(id, existingBlog)
	return existingBlog, err
}

// IMPLEMENTACIJA NOVE METODE
func (s *blogService) UpdateComment(blogID, commentID primitive.ObjectID, updatedComment *domain.Comment, userID string) error {
	// 1. Dobavljamo blog da pronađemo komentar
	blog, err := s.repo.GetById(blogID)
	if err != nil {
		return err
	}

	// 2. Pronalazimo komentar i proveravamo autorstvo
	var targetComment *domain.Comment
	for i := range blog.Comments {
		if blog.Comments[i].ID == commentID {
			targetComment = &blog.Comments[i]
			break
		}
	}

	if targetComment == nil {
		return errors.New("comment not found")
	}

	if targetComment.AuthorID != userID {
		return errors.New("forbidden")
	}

	// 3. Ako je sve u redu, ažuriramo podatke
	targetComment.Text = updatedComment.Text
	targetComment.LastUpdatedAt = time.Now()

	return s.repo.UpdateComment(blogID, targetComment)
}

// IMPLEMENTACIJA NOVE METODE ZA SAGA OBRAZAC
func (s *blogService) RemoveLikesFromAuthorBlogs(userID, authorID string) error {
	// 1. Dobavljamo sve blogove određenog autora
	blogs, err := s.repo.GetBlogsByAuthors([]string{authorID})
	if err != nil {
		return err
	}

	// 2. Za svaki blog uklanjamo lajkove određenog korisnika
	for _, blog := range blogs {
		// Proveravamo da li korisnik ima lajk na ovom blogu
		hasLike := false
		for _, like := range blog.Likes {
			if like.UserID == userID {
				hasLike = true
				break
			}
		}

		// Ako ima lajk, uklanjamo ga
		if hasLike {
			err := s.repo.RemoveLike(blog.ID, userID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
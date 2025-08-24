// u services/blog-service/repository/blog_repository.go

package repository

import (
	"context"
	"log" // <-- OVAJ IMPORT JE NEDOSTAJAO


	"blog-service/domain" // Importujemo naš model

	"go.mongodb.org/mongo-driver/bson" // <-- I OVAJ IMPORT JE NEDOSTAJAO
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// Definišemo koje sve operacije repozitorijum može da izvrši
type BlogRepository interface {
	Create(blog *domain.Blog) error
	GetById(id primitive.ObjectID) (*domain.Blog, error)
	GetAll() ([]*domain.Blog, error)
	AddComment(blogID primitive.ObjectID, comment *domain.Comment) error // <-- NOVA METODA
	AddLike(blogID primitive.ObjectID, like *domain.Like) error       // <-- NOVA METODA
	RemoveLike(blogID primitive.ObjectID, userID string) error // <-- NOVA METODA
	Update(id primitive.ObjectID, blog *domain.Blog) error // <-- NOVA METODA
	UpdateComment(blogID primitive.ObjectID, comment *domain.Comment) error // <-- NOVA METODA
}

// Struktura koja implementira interfejs
type blogRepository struct {
	blogs *mongo.Collection
}

// Konstruktor koji kreira novu instancu repozitorijuma
func NewBlogRepository(client *mongo.Client) BlogRepository {
	// Pristupamo bazi 'soa_db' i kolekciji 'blogs'
	// Ako ne postoje, Mongo će ih automatski kreirati
	blogs := client.Database("soa_db").Collection("blogs")
	return &blogRepository{blogs: blogs}
}

// Implementacija prve metode - kreiranje novog bloga
func (r *blogRepository) Create(blog *domain.Blog) error {
	// Umećemo jedan dokument (blog) u kolekciju
	_, err := r.blogs.InsertOne(context.TODO(), blog)
	return err
}

// Za sada ostavljamo ostale metode neimplementirane
func (r *blogRepository) GetById(id primitive.ObjectID) (*domain.Blog, error) {
	filter := bson.M{"_id": id}
	var blog domain.Blog
	err := r.blogs.FindOne(context.TODO(), filter).Decode(&blog)
	
	// Eksplicitno proveravamo da li je greška "not found" i vraćamo je
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err // Vraćamo specifičnu grešku
		}
		// Vraćamo bilo koju drugu grešku
		return nil, err
	}
	
	return &blog, nil
}

func (r *blogRepository) GetAll() ([]*domain.Blog, error) {
	cursor, err := r.blogs.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error finding blogs:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []*domain.Blog
	// Iteriramo kroz sve rezultate iz baze
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		log.Println("Error decoding blogs:", err)
		return nil, err
	}

	return blogs, nil
}

// IMPLEMENTACIJA NOVE METODE
func (r *blogRepository) AddComment(blogID primitive.ObjectID, comment *domain.Comment) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$push": bson.M{"comments": comment}}

	_, err := r.blogs.UpdateOne(context.TODO(), filter, update)
	return err
}

// IMPLEMENTACIJA AddLike
func (r *blogRepository) AddLike(blogID primitive.ObjectID, like *domain.Like) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$push": bson.M{"likes": like}}
	_, err := r.blogs.UpdateOne(context.TODO(), filter, update)
	return err
}

// IMPLEMENTACIJA RemoveLike
func (r *blogRepository) RemoveLike(blogID primitive.ObjectID, userID string) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$pull": bson.M{"likes": bson.M{"userId": userID}}}
	_, err := r.blogs.UpdateOne(context.TODO(), filter, update)
	return err
}

// IMPLEMENTACIJA NOVE METODE
func (r *blogRepository) Update(id primitive.ObjectID, blog *domain.Blog) error {
	filter := bson.M{"_id": id}
	// Ažuriramo samo naslov i sadržaj, jer samo to autor može da menja
	update := bson.M{"$set": bson.M{
		"title":   blog.Title,
		"content": blog.Content,
		"images":  blog.Images,
	}}
	_, err := r.blogs.UpdateOne(context.TODO(), filter, update)
	return err
}

// IMPLEMENTACIJA NOVE METODE
func (r *blogRepository) UpdateComment(blogID primitive.ObjectID, comment *domain.Comment) error {
	filter := bson.M{"_id": blogID}
	update := bson.M{"$set": bson.M{
		"comments.$[elem].text":          comment.Text,
		"comments.$[elem].lastUpdatedAt": comment.LastUpdatedAt,
	}}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": comment.ID}},
	})

	_, err := r.blogs.UpdateOne(context.TODO(), filter, update, arrayFilters)
	return err
}
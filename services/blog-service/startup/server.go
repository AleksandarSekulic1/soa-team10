package startup

import (
	"blog-service/api"
	// 1. Dodajemo import za docs (koji će uskoro biti kreiran)
	_ "blog-service/docs"
	"blog-service/repository"
	"blog-service/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// 2. Dodajemo importe za swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	mongoClient := ConnectDB()

	blogRepo := repository.NewBlogRepository(mongoClient)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := api.NewBlogHandler(blogService)

	router := gin.Default()
	
	// CORS ... (ostaje isto)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// 3. Dodajemo rutu za Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

apiRoutes := router.Group("/api/blogs")
    {
        // Primenjujemo middleware na POST rutu
        apiRoutes.POST("", api.AuthMiddleware(), blogHandler.CreateBlog)
        // GET ruta je sada takođe zaštićena
        apiRoutes.GET("", api.AuthMiddleware(), blogHandler.GetAllBlogs)
        // Nova ruta za blogove od praćenih korisnika
        apiRoutes.GET("/following", api.AuthMiddleware(), blogHandler.GetBlogsFromFollowing)
        apiRoutes.POST("/:id/comments", api.AuthMiddleware(), blogHandler.AddComment)
        apiRoutes.POST("/:id/likes", api.AuthMiddleware(), blogHandler.ToggleLike)
        // GET jednog bloga je sada takođe zaštićena ruta
        apiRoutes.GET("/:id", api.AuthMiddleware(), blogHandler.GetBlogById) 
        apiRoutes.PUT("/:id", api.AuthMiddleware(), blogHandler.UpdateBlog)
        apiRoutes.PUT("/:id/comments/:commentId", api.AuthMiddleware(), blogHandler.UpdateComment)
    }

	return &Server{router: router}
}

func (s *Server) Start() {
	err := s.router.Run(":8082")
	if err != nil {
		panic(err)
	}
}
package startup

import (
	"tours-service/api"
	_ "tours-service/docs"
	"tours-service/repository"
	"tours-service/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	// 1. Povezujemo se na bazu na početku
	mongoClient := ConnectDB()

	// 2. Inicijalizujemo sve slojeve, kao u vašem primeru
	tourRepo := repository.NewTourRepository(mongoClient)
	tourService := service.NewTourService(tourRepo)
	tourHandler := api.NewTourHandler(tourService)

	// 3. Kreiramo ruter i podešavamo CORS
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Dozvoljavamo sve za sada, možete promeniti na http://localhost:4200
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// 4. Definišemo Swagger i API rute
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiGroup := router.Group("/api")
	{
		toursGroup := apiGroup.Group("/tours")
		{
			toursGroup.POST("", api.AuthMiddleware(), tourHandler.Create)
			toursGroup.GET("/my-tours", api.AuthMiddleware(), tourHandler.GetByAuthor)
			toursGroup.GET("", tourHandler.GetAll) // Javna ruta za prikaz svih tura
			toursGroup.POST("/:id/reviews", api.AuthMiddleware(), tourHandler.AddReview)
		}
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	err := s.router.Run(":8083")
	if err != nil {
		panic(err)
	}
}

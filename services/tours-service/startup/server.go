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

    //touristPositionRepo := repository.NewTouristPositionRepository(mongoClient)
	//touristPositionService := service.NewTouristPositionService(touristPositionRepo)
	//touristPositionHandler := api.NewTouristPositionHandler(touristPositionService)

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
			toursGroup.GET("/:id", tourHandler.GetById)
			// NOVE RUTE ZA KEY POINTS
			toursGroup.POST("/:id/keypoints", api.AuthMiddleware(), tourHandler.AddKeyPoint)
			toursGroup.PUT("/:id/keypoints/:keypointId", api.AuthMiddleware(), tourHandler.UpdateKeyPoint)
			toursGroup.DELETE("/:id/keypoints/:keypointId", api.AuthMiddleware(), tourHandler.DeleteKeyPoint)
			// --- NOVE RUTE ZA STANJA TURE ---
			toursGroup.POST("/:id/transport-info", api.AuthMiddleware(), tourHandler.AddTransportInfo)
			toursGroup.POST("/:id/publish", api.AuthMiddleware(), tourHandler.Publish)
			toursGroup.POST("/:id/archive", api.AuthMiddleware(), tourHandler.Archive)
			toursGroup.GET("/published", tourHandler.GetPublished)
			toursGroup.POST("/:id/reactivate", api.AuthMiddleware(), tourHandler.Reactivate)
			toursGroup.GET("/archived", tourHandler.GetArchived) // NOVO: Registracija rute
		}

		/*positionGroup := apiGroup.Group("/tourist-position")
		{
			// Obe rute zahtevaju da je korisnik ulogovan
			positionGroup.GET("", api.AuthMiddleware(), touristPositionHandler.GetByUserId)
			positionGroup.POST("", api.AuthMiddleware(), touristPositionHandler.Update)
		}*/
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	err := s.router.Run(":8083")
	if err != nil {
		panic(err)
	}
}

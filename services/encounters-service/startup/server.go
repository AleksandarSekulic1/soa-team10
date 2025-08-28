package startup

import (
	"encounters-service/api"
	"encounters-service/repository"
	"encounters-service/service"
	
	_ "encounters-service/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	mongoClient := ConnectDB()

	// Inicijalizacija za Simulator Pozicije
	touristPositionRepo := repository.NewTouristPositionRepository(mongoClient)
	touristPositionService := service.NewTouristPositionService(touristPositionRepo)
	touristPositionHandler := api.NewTouristPositionHandler(touristPositionService)
	
	// --- NOVO: Inicijalizacija za Izvođenje Ture ---
	tourExecutionRepo := repository.NewTourExecutionRepository(mongoClient)
	tourExecutionService := service.NewTourExecutionService(tourExecutionRepo)
	tourExecutionHandler := api.NewTourExecutionHandler(tourExecutionService)
	
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := router.Group("/api")
	{
		// Rute za Simulator Pozicije
		positionGroup := apiGroup.Group("/tourist-position")
		{
			positionGroup.GET("", api.AuthMiddleware(), touristPositionHandler.GetByUserId)
			positionGroup.POST("", api.AuthMiddleware(), touristPositionHandler.Update)
		}

		// --- NOVO: Rute za Izvođenje Ture ---
		executionGroup := apiGroup.Group("/tour-executions")
		{
			// Sve rute zahtevaju da je korisnik ulogovan
			executionGroup.POST("/start/:tourId", api.AuthMiddleware(), tourExecutionHandler.StartTour)
			executionGroup.POST("/check-position", api.AuthMiddleware(), tourExecutionHandler.CheckPosition)
			executionGroup.POST("/:executionId/complete", api.AuthMiddleware(), tourExecutionHandler.CompleteTour)
			executionGroup.POST("/:executionId/abandon", api.AuthMiddleware(), tourExecutionHandler.AbandonTour)
		}
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	err := s.router.Run(":8084") 
	if err != nil {
		panic(err)
	}
}
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
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router *gin.Engine
}

func NewServer(client *mongo.Client) *Server {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tourRepo := repository.NewTourRepository(client)
	tourService := service.NewTourService(tourRepo)
	tourHandler := api.NewTourHandler(tourService)

	apiGroup := router.Group("/api")
	{
		toursGroup := apiGroup.Group("/tours")
		{
			toursGroup.POST("", api.AuthMiddleware(), tourHandler.Create)
		}
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	s.router.Run(":8082")
}

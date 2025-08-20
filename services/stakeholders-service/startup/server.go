// startup/server.go

package startup

import (
	"stakeholders-service/api"
	_ "stakeholders-service/docs"
	"stakeholders-service/repository"
	"stakeholders-service/service"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j" // Import drajvera
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

// NewServer sada prihvata drajver
func NewServer(driver neo4j.DriverWithContext) *Server {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prosleđujemo drajver u repozitorijum
	userRepo := repository.NewUserRepository(driver)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	apiGroup := router.Group("/api")
	{
		stakeholdersGroup := apiGroup.Group("/stakeholders")
		{
			stakeholdersGroup.POST("/register", userHandler.Register)
		}
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	s.router.Run(":8081")
}

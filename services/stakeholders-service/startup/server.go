// startup/server.go

package startup

import (
	"stakeholders-service/api"
	_ "stakeholders-service/docs"
	"stakeholders-service/repository"
	"stakeholders-service/service"

	"github.com/gin-contrib/cors" // Uverite se da je ovaj import tu
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer(driver neo4j.DriverWithContext) *Server {
	router := gin.Default()

	// Detaljna CORS konfiguracija
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Dozvoli zahteve sa Angular aplikacije
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // EKSPLICITNO DOZVOLI AUTHORIZATION HEADER

	router.Use(cors.New(config)) // Koristimo novu, detaljnu konfiguraciju

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := repository.NewUserRepository(driver)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	apiGroup := router.Group("/api")
	{
		stakeholdersGroup := apiGroup.Group("/stakeholders")
		{
			stakeholdersGroup.POST("/register", userHandler.Register)
			stakeholdersGroup.GET("", userHandler.GetAll)
			stakeholdersGroup.POST("/login", userHandler.Login)
		}
	}

	return &Server{router: router}
}

func (s *Server) Start() {
	s.router.Run(":8081")
}

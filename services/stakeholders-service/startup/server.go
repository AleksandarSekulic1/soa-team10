package startup

import (
	"fmt"
	"stakeholders-service/api"
	_ "stakeholders-service/docs"
	"stakeholders-service/repository"
	"stakeholders-service/service"

	"github.com/gin-contrib/cors"
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
	router.RedirectTrailingSlash = false

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := repository.NewUserRepository(driver)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	// EKSPLICITNA I PREGLEDNA STRUKTURA RUTA
	apiRoutes := router.Group("/api")
	{
		// Javne rute
		apiRoutes.POST("/stakeholders/register", userHandler.Register)
		apiRoutes.POST("/stakeholders/login", userHandler.Login)

		// Rute za ulogovane korisnike
		apiRoutes.GET("/stakeholders/profile", api.AuthMiddleware(), userHandler.GetProfile)
		apiRoutes.PUT("/stakeholders/profile", api.AuthMiddleware(), userHandler.UpdateProfile)

		// Rute samo za administratore
		apiRoutes.GET("/stakeholders", api.AuthMiddleware(), api.AdminRoleMiddleware(), userHandler.GetAll)
		apiRoutes.PUT("/stakeholders/:username/block", api.AuthMiddleware(), api.AdminRoleMiddleware(), userHandler.BlockUser)
		apiRoutes.PUT("/stakeholders/:username/unblock", api.AuthMiddleware(), api.AdminRoleMiddleware(), userHandler.UnblockUser)
	}

	// Ostavljamo ispis ruta radi provere
	fmt.Println("--- REGISTROVANE RUTE ---")
	for _, route := range router.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	fmt.Println("-------------------------")

	return &Server{router: router}
}

func (s *Server) Start() {
	s.router.Run(":8081")
}

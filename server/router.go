package server

import (
	"fmt"
	"os"
	"time"

	"github.com/fabianpoels/fabianpoels-api-go/controllers"
	"github.com/fabianpoels/fabianpoels-api-go/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// FOR DEV PURPOSES!!!!
	// TODO: rework to dynamically configure cors, depending on the env
	corsConfig := cors.DefaultConfig()

	if os.Getenv("environment") == "development" {
		// LOCAL DEV CONFIG
		// domain := config.GetEnv("DOMAIN")
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{"http://localhost:9000", "http://127.0.0.1:9000"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "PATCH", "GET", "PUT", "DELETE", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	} else {
		domain := os.Getenv("DOMAIN")
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{fmt.Sprintf("http://%s", domain)}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "PATCH", "GET", "PUT", "DELETE", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	}

	router.Use(cors.New(corsConfig))

	// CONTROLLERS
	authController := new(controllers.AuthController)

	// ROUTING
	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.POST("/auth/login", authController.Login)
			v1.POST("/auth/refresh-token", authController.RefreshToken)
			v1.POST("/auth/logout", middleware.ValidateJwt(), authController.Logout)

			// LOGGED IN
			v1.Use(middleware.ValidateJwt())
			{

			}
		}
	}
	return router

}

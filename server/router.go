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
	publicController := new(controllers.PublicController)
	adminController := new(controllers.AdminController)

	// ROUTING
	api := router.Group("api")
	{
		// AUTH
		api.POST("/auth/login", authController.Login)
		api.POST("/auth/refresh-token", authController.RefreshToken)
		api.POST("/auth/logout", middleware.ValidateJwt(), authController.Logout)

		// PUBLIC
		public := api.Group("public")
		{
			public.GET("ascents", publicController.Ascents)
		}

		// LOGGED IN
		api.Use(middleware.ValidateJwt())
		{
			admin := api.Group("admin")
			{
				admin.POST("ascent", adminController.AddAscent)
			}
		}

	}
	return router

}

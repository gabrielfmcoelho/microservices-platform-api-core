package main

import (
	"log"
	"time"

	"github.com/gabrielfmcoelho/platform-core/api/route"
	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           Platform API
// @version         0.1.1
// @description		Platform API is a RESTful API for managing ...
// @termsOfService  http://swagger.io/terms/

// @contact.name   Eng. Gabriel Coelho | InovaData development team
// @contact.url    https://solude.tech
// @contact.email  suporte@solude.tech

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Initialize the application
	app := bootstrap.App()
	defer app.CloseDBConnection()

	// Configuration variables
	env := app.Env

	// Database instance (Gorm DB)
	db := app.DB

	// Context timeout
	timeout := time.Duration(env.ContextTimeout) * time.Second

	// Create a Gin router instancehttps://github.com/inova-data-tech/Solude-api.git
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://195.200.0.244:3000", "http://solude.inovadata.tech", "https://solude.inovadata.tech", "https://solude.tech"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Route binding
	route.Setup(env, timeout, db, router)

	// Run the server
	if err := router.Run(env.ServerAddress); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

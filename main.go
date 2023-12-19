package main

import (
	"log"
	"net/http"
	"os"
	"profile-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Update with your allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	// Set up Logger middleware
	router.Use(gin.Logger())

	// Initialize routes
	routes.InitRoutes(router)

	// Set up server port
	port := os.Getenv("PORT")

	// Start the server
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}

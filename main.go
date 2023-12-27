package main

import (
	"log"
	"net/http"
	"os"
	"profile-api/middleware"
	"profile-api/routes"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

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

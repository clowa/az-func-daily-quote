package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIPort() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func sayHello(c *gin.Context) {
	log.Println("Invoke ROOT")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, from Azure",
	})
	// return
}

func about(c *gin.Context) {
	log.Println("Invoke ABOUT")
	routes := [...]string{"/", "/about", "/health"}
	c.JSON(http.StatusOK, gin.H{
		"message": "Exported routes",
		"routes":  routes,
	},
	)
	// return
}

func health(c *gin.Context) {
	log.Println("Invoke HEALTH")
	c.JSON(http.StatusOK, gin.H{
		"message": "Everything is fantastic with this service",
	})
}

func main() {
	router := gin.Default()
	router.GET("api/simpleapi", sayHello)
	router.GET("api/simpleapi/health", health)
	router.GET("api/simpleapi/about", about)

	port_info := APIPort()
	router.Run(port_info)
	log.Println("API is up & running - " + port_info)
}

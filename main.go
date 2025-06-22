package main

import (
	"backend-ventas/db"
	"backend-ventas/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Corriendo en localhost:8080 !!!")

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	db.Conectar() // conecta la base de datos
	routes.ClienteRoutes(router)

	router.Run(":7777") // listen and serve on 0.0.0.0:8080
}

package main

import (
	"backend-ventas/db"
	"backend-ventas/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	db.Conectar() // conecta la base de datos desde la funcion Conectar() del archivo db.go
	routes.ClienteRoutes(router) 

	puerto := ":8080"
	router.Run(puerto) // listen and serve on 0.0.0.0:8080
	fmt.Print("\n\n\t\t>>>>> Corriendo en localhost:", puerto ,"!!!\n\n")
}

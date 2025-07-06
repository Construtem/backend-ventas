package main

import (
	"backend-ventas/api/database"
	//"backend-ventas/api/routes"

	apiRoutes "backend-ventas/api/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Conectar a la base de datos
	database.InitDB()

	// Configurar todas las rutas usando SetupRoutes
	apiRoutes.SetupRoutes(router)
	apiRoutes.ClienteRoutes(router) 

	puerto := ":8080"
	router.Run(puerto) // listen and serve on 0.0.0.0:8080
	fmt.Print("\n\n\t\t>>>>> Corriendo en localhost:", puerto, "!!!\n\n")
}

package main

import (
	//"simplerestapi/controllers"

	"github.com/gin-gonic/gin"

	"simplerestapi/routes"

	"simplerestapi/database"
)

func main() {	
   r := setupRouter()
	r.Run(":5000")
}

func setupRouter() *gin.Engine {
	db := database.InitDb()
   r := gin.Default()
	routes.SetGetPong(r)
	routes.SetUserRoute(r,db)
	routes.SetAuthRoute(r,db)

	return r
}
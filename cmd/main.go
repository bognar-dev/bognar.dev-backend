package main

import (
	"bognar.dev-backend/controllers"
	"bognar.dev-backend/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	database.InitDB()
	r := gin.Default()

	// Serve JSON data
	r.GET("/api/projects", controllers.GetProjects)

	// Create a JSON file in the "projects" folder
	r.POST("/api/createProject", controllers.CreateProject)

	r.Run()
}

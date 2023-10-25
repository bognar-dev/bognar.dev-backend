package main

import (
	"bognar.dev-backend/controllers"
	"bognar.dev-backend/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	database.InitDB()
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()
	var secret = []byte("secret")
	// Setup the cookie store for session management
	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	// Login and logout, register routes
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.POST("/register", controllers.Register)

	// Serve projects
	r.GET("/api/projects", controllers.GetProjects)

	// Private group, require authentication to access
	private := r.Group("/private")
	private.Use(controllers.AuthRequired)
	{
		private.GET("/me", controllers.Me)
		private.GET("/status", controllers.Status)
		private.POST("/createProject", controllers.CreateProject)
	}
	return r
}

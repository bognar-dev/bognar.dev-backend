package main

import (
	"bognar.dev-backend/controllers"
	"bognar.dev-backend/database"
	"bognar.dev-backend/middlewares"
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
	gin.SetMode(gin.ReleaseMode)
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()
	var secret = []byte("secret")
	// Setup the cookie store for session management
	//r.Use(middlewares.RateLimit)
	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	// Login and logout, register routes
	r.GET("/", controllers.Hey)
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	// Serve projects
	r.GET("/api/projects", controllers.GetProjects)

	// Private group, require authentication to access
	private := r.Group("/private")
	private.Use(middlewares.JwtAuthMiddleware())
	{
		private.GET("/user", controllers.CurrentUser)
		private.GET("/status", controllers.Status)
		private.POST("/createProject", controllers.CreateProject)
	}
	return r
}

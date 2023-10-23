package controllers

import (
	"bognar.dev-backend/database"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectData struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Url             string   `json:"url"`
	LongDescription string   `json:"longDescription"`
	Tags            []string `json:"tags"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Status          string   `json:"status"`
	TeamMembers     []string `json:"teamMembers"`
	GithubRepo      string   `json:"githubRepo"`
	Image           string   `json:"image"`
}
type Project struct {
	ID        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	Data      string `db:"data" json:"data"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

func GetProjects(c *gin.Context) {
	var project []Project
	err := database.DBClient.Select(&project, "SELECT data::json FROM projects ")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var projects []ProjectData
	for _, p := range project {
		var projectDetails ProjectData
		err = json.Unmarshal([]byte(p.Data), &projectDetails)
		projects = append(projects, projectDetails)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, projects)
}
func CreateProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectJSON, err := json.Marshal(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	_, err = database.DBClient.Exec(`INSERT INTO projects (data)
        VALUES ($1)`, projectJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Print(projectJSON)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

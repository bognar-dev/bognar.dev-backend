package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

type Project struct {
	Id              int      `json:"id"`
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

func readProjectsFromDirectory(directoryPath string) ([]Project, error) {
	var projectList []Project

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			project, err := readProjectFile(path)
			if err != nil {
				return err
			}
			projectList = append(projectList, project)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return projectList, nil
}
func readProjectFile(filePath string) (Project, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Project{}, err
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return Project{}, err
	}

	return project, nil
}

func main() {
	r := gin.Default()

	// Serve JSON data
	r.GET("/api/projects", func(c *gin.Context) {
		projectList, err := readProjectsFromDirectory("./projects")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, projectList)
	})

	// Create a JSON file in the "projects" folder
	r.POST("/api/createProject", func(c *gin.Context) {
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

		projectDir := "./projects"
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			if err := os.Mkdir(projectDir, 0755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
				return
			}
		}

		filePath := fmt.Sprintf("%s/%d.json", projectDir, project.Id)
		err = os.WriteFile(filePath, projectJSON, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JSON file"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
	})

	r.Run(":8080")
}

package controllers

import (
	"bognar.dev-backend/database"
	"bognar.dev-backend/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project
	err := database.DBClient.Select(&projects, "SELECT * FROM projects ")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*var projectsJSONBlob []models.ProjectData
	for _, p := range projects {
		var projectDetails models.ProjectData
		err = json.Unmarshal([]byte(p.Data), &projectDetails)
		projectsJSONBlob = append(projectsJSONBlob, projectDetails)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}*/
	fmt.Println(projects)

	c.JSON(http.StatusOK, projects)
}
func CreateProject(c *gin.Context) {
	var project models.ProjectData
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectJSON, err := json.Marshal(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}
	fmt.Println(project)

	_, err = database.DBClient.Exec(`INSERT INTO projects (data)
        VALUES ($1)`, projectJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Print(projectJSON)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}
func GetProjectByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	var project models.Project
	err := database.DBClient.Get(&project, "SELECT * FROM projects WHERE id = ($1)", id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}
func UpdateProject(c *gin.Context) {

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(project)

	err, _ := database.DBClient.NamedExec(`UPDATE projects SET
                    data=:data,
                    updated_at=: time.Now()
                WHERE id = :ID`, project)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, project)
}

func Hey(c *gin.Context) {
	fmt.Println("hey")
	c.JSON(http.StatusOK, gin.H{"data": "Hey!!!"})
}

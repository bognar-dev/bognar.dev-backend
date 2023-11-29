package controllers

import (
	"bognar.dev-backend/database"
	"bognar.dev-backend/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	supa "github.com/nedpals/supabase-go"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	fmt.Println("Hello from UpdateProject")
	var updateForm models.UpdateProjectForm
	var project models.Project
	err := c.Bind(&updateForm)
	if err != nil {
		fmt.Println("Form Bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	filePath := filepath.Join("uploads", updateForm.Image.Filename)
	// Save the file to the defined path
	if err := c.SaveUploadedFile(updateForm.Image, filePath); err != nil {
		fmt.Println("SaveUploadedFile error:", err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	storage := database.SBClient.Storage.From("images")
	fmt.Println("storage:", storage.BucketId)
	var update = false
	files := storage.List("", supa.FileSearchOptions{})
	for _, f := range files {
		if f.Name == updateForm.Image.Filename {
			fmt.Println("File already exists")
			update = true
		}
	}
	updateConfirm := storage.UploadOrUpdate("/"+updateForm.Image.Filename, file, update)
	fmt.Println("Upload Message = ", updateConfirm.Message, " Upload Key = ", updateConfirm.Key)
	publicURL := storage.GetPublicUrl(updateConfirm.Key)
	fmt.Println("Public URL = ", publicURL)

	project.Data = models.ProjectData{
		Image:           publicURL.SignedUrl,
		Name:            updateForm.ProjectName,
		Description:     updateForm.Description,
		LongDescription: updateForm.LongDescription,
		StartDate:       updateForm.SinceDate,
		Tags:            updateForm.Tags,
	}

	project.Id, _ = strconv.Atoi(updateForm.ID)
	project.UpdatedAt = time.Now()

	res, err := database.DBClient.NamedExec(`UPDATE projects SET
		                    data=:data,
		                    updated_at=:updated_at
		                WHERE id = :id`, &project)
	fmt.Println("Update res:", res)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(project)
	c.JSON(http.StatusOK, project)
}

func Hey(c *gin.Context) {
	fmt.Println("hey")
	c.JSON(http.StatusOK, gin.H{"data": "Hey!!!"})
}

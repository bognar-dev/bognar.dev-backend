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
	"strings"
	"time"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project
	err := database.DBClient.Select(&projects, "SELECT * FROM projects ")
	fmt.Println("Projects ", projects)
	fmt.Println("err ", err)
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
	var filename string
	filename = strings.ToLower(updateForm.Image.Filename)
	filename = strings.ReplaceAll(filename, " ", "-")
	fmt.Println("storage:", storage.BucketId)
	var update = false
	files := storage.List("", supa.FileSearchOptions{})
	for _, f := range files {
		if f.Name == filename {
			fmt.Println("File already exists")
			update = true
		}
	}
	var signedUrl string
	if updateForm.Image.Filename != "undefined" {
		updateConfirm := storage.UploadOrUpdate("/"+filename, file, update)
		fmt.Println("Upload Message = ", updateConfirm.Message, " Upload Key = ", updateConfirm.Key)
		if updateConfirm.Message != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}
		UploadKey := strings.Split(updateConfirm.Key, "/")[1]
		fmt.Println("UploadKey:", UploadKey)
		publicURL := storage.GetPublicUrl(UploadKey)
		signedUrl = publicURL.SignedUrl
		fmt.Println("Public URL = ", publicURL)
		fmt.Println("updateFormTags:", updateForm.Tags)
		fmt.Println("Signed URL:", publicURL.SignedUrl)
	} else {
		signedUrl = updateForm.ImageURL
	}
	var tagArray []string

	err = json.Unmarshal([]byte(updateForm.Tags), &tagArray)
	if err != nil {
		fmt.Println(err)
		return
	}

	project.Data = models.ProjectData{
		Image:           signedUrl,
		Name:            updateForm.ProjectName,
		Description:     updateForm.Description,
		LongDescription: updateForm.LongDescription,
		StartDate:       updateForm.SinceDate,
		Tags:            tagArray,
	}

	project.Id, _ = strconv.Atoi(updateForm.ID)
	project.UpdatedAt = time.Now()
	fmt.Println("Project:", project.Id)
	res, err := database.DBClient.NamedExec(`UPDATE projects SET
		                    data=:data,
		                    updated_at=:updated_at
		                WHERE id = :id`, &project)
	var row, _ = res.RowsAffected()
	fmt.Println("Update res:", row)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(project)
	c.JSON(http.StatusOK, "Project updated successfully")
}

func Hey(c *gin.Context) {
	fmt.Println("hey")
	c.JSON(http.StatusOK, gin.H{"data": "Hey!!!"})
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"mime/multipart"
	"time"
)

type UpdateProjectForm struct {
	ID              string                `form:"id"`
	ImageURL        string                `form:"imageURL"`
	Image           *multipart.FileHeader ` form:"image"`
	ProjectName     string                `form:"projectName"`
	SinceDate       string                `form:"sinceDate"`
	Tags            []string              `form:"tags"`
	Description     string                `form:"description"`
	LongDescription string                `form:"longDescription"`
}

type ProjectData struct {
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
	Id        int         `db:"id" json:"id"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	Data      ProjectData `db:"data" json:"data"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

func (p *Project) UnmarshalJSON(data []byte) error {
	type Alias Project // Create an alias to avoid recursion

	aux := &struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.CreatedAt != "" {
		createdTime, err := time.Parse("2006-01-02 15:04:05.999999-07", aux.CreatedAt)

		if err != nil {
			return err
		}
		p.CreatedAt = createdTime
	}

	if aux.UpdatedAt != "" {
		updatedTime, err := time.Parse("2023-10-27 16:53:56.565619+00", aux.UpdatedAt)
		if err != nil {
			return err
		}
		p.UpdatedAt = updatedTime
	}

	return nil
}

// Implement the Value method to convert ProjectData to a database value.
func (pd ProjectData) Value() (driver.Value, error) {
	// Marshal the ProjectData as JSON
	jsonData, err := json.Marshal(pd)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// Implement the Scan method to convert a database value to ProjectData.
func (pd *ProjectData) Scan(value interface{}) error {
	// Ensure the value is a byte slice
	byteData, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}

	// Unmarshal JSON into ProjectData
	if err := json.Unmarshal(byteData, pd); err != nil {
		return err
	}
	return nil
}

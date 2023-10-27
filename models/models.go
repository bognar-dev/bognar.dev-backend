package models

import (
	"bognar.dev-backend/database"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

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
	ID        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	Data      string `db:"data" json:"data"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID         int    `db:"id" json:"id"`
	SignedUpAt string `db:"signed_up_at" json:"signed_up_at"`
	Username   string `db:"username" json:"username"`
	Password   string `db:"password" json:"password"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) error {

	var err error

	var u User

	err = database.DBClient.Get(
		&u,
		"SELECT * FROM users LIMIT 1")
	fmt.Print(err)
	if err != nil {
		return err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	return nil

}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = u.BeforeSave()

	if err != nil {
		return nil, err
	}
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = ($1)"

	err = database.DBClient.Get(&count, query, u.Username)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("USer Already exists")
	}
	_, err = database.DBClient.NamedExec(`INSERT INTO users (username, password) VALUES (:username, :password)`, u)
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

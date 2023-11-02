package models

import (
	"bognar.dev-backend/database"
	"bognar.dev-backend/utils"
	"database/sql/driver"
	"encoding/json"
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
	ID        int         `db:"id" json:"id"`
	CreatedAt string      `db:"created_at" json:"created_at"`
	Data      ProjectData `db:"data" json:"data"`
	UpdatedAt string      `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID         uint   `db:"id" json:"id"`
	SignedUpAt string `db:"signed_up_at" json:"signed_up_at"`
	Username   string `db:"username" json:"username"`
	Password   string `db:"password" json:"password"`
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

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(user *User) (string, error) {

	var err error

	var databaseUser User
	//var preparedUsername = "'"+user.Username+"'"
	//fmt.Println()
	err = database.DBClient.Get(
		&databaseUser,
		`SELECT * FROM users WHERE username = $1`, user.Username)
	fmt.Println("databaseuser ", databaseUser)
	fmt.Println("database error = ", err)
	if err != nil {
		return "", err
	}

	err = VerifyPassword(user.Password, databaseUser.Password)
	fmt.Println("password verified")
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", err
	}

	token, err := token.GenerateToken(databaseUser.ID)

	if err != nil {
		return "", err
	}
	fmt.Println("token generated")
	return token, nil

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

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := database.DBClient.Get(&u, `SELECT * from users where id = ($1)`, uid); err != nil {
		return u, err
	}

	u.PrepareGive()

	return u, nil

}
func (u *User) PrepareGive() {
	u.Password = ""
}

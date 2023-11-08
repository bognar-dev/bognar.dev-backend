package models

import (
	"bognar.dev-backend/database"
	"bognar.dev-backend/utils"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

type User struct {
	ID         uint   `db:"id" json:"id"`
	SignedUpAt string `db:"signed_up_at" json:"signed_up_at"`
	Username   string `db:"username" json:"username"`
	Password   string `db:"password" json:"password"`
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

	t, err := token.GenerateToken(databaseUser.ID)

	if err != nil {
		return "", err
	}
	fmt.Println("t generated")
	return t, nil

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

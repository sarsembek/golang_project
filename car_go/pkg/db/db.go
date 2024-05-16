package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
    "car_project/pkg/model"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "postgres://postgres:postgres@postgres/cars?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// CreateUser inserts a new user into the database
func CreateUser(user model.User) error {
	_, err := DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

// GetUserByUsername retrieves a user by username from the database
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AuthenticateUser checks if the given login credentials are valid
func AuthenticateUser(username, password string) (bool, *model.User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false, nil, err
	}

	isAuthenticated, err := user.Authenticate(password)
	if err != nil {
		return false, nil, err
	}

	return isAuthenticated, user, nil
}

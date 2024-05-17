package db

import (
	"database/sql"
	"log"

	"car_project/pkg/model"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
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
	// Run migrations
	if err := runMigrations(DB); err != nil {
		log.Fatalf("could not apply migrations: %v", err)
	}

}

//Run migrations
func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://path/to/migrations", // Replace with the path to your migration files
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
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

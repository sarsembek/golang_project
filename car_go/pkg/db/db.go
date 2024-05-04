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
	DB, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/cars?sslmode=disable")
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

func CreateCar(c model.Car) error {
	_, err := DB.Exec("INSERT INTO cars (brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
			c.Brand,
			c.Model,
			c.Year,
			c.Color,
			c.BodyStyle,
			c.EngineSize,
			c.Weight,
			c.BasePrice,
			c.FuelCapacity,
			c.Horsepower,
			c.Torque,
			c.Acceleration,
			c.TopSpeed,
	)
	if err != nil {
			return err
	}
	return nil
}

func GetAllCars() ([]model.Car, error) {
	var cars []model.Car
	rows, err := DB.Query("SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars")
	if err != nil {
			return cars, err
	}
	defer rows.Close()

	for rows.Next() {
			var c model.Car
			if err := rows.Scan(
					&c.ID,
					&c.Brand,
					&c.Model,
					&c.Year,
					&c.Color,
					&c.BodyStyle,
					&c.EngineSize,
					&c.Weight,
					&c.BasePrice,
					&c.FuelCapacity,
					&c.Horsepower,
					&c.Torque,
					&c.Acceleration,
					&c.TopSpeed,
			); err != nil {
					return cars, err
			}
			cars = append(cars, c)
	}
	if err := rows.Err(); err != nil {
			return cars, err
	}
	return cars, nil
}

func GetCarWithPagination(page, limit int, sortBy, filterBy string) ([]model.Car, error) {
	var cars []model.Car

	query := "SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars"

	if filterBy != "" {
			query += " WHERE brand LIKE '%" + filterBy + "%' OR model LIKE '%" + filterBy + "%'"
	}

	if sortBy != "" {
			query += " ORDER BY " + sortBy
	}

	query += " LIMIT $1 OFFSET $2"
	offset := (page - 1) * limit

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
			return cars, err
	}
	defer rows.Close()

	for rows.Next() {
			var c model.Car
			if err := rows.Scan(
					&c.ID,
					&c.Brand,
					&c.Model,
					&c.Year,
					&c.Color,
					&c.BodyStyle,
					&c.EngineSize,
					&c.Weight,
					&c.BasePrice,
					&c.FuelCapacity,
					&c.Horsepower,
					&c.Torque,
					&c.Acceleration,
					&c.TopSpeed,
			); err != nil {
					return cars, err
			}
			cars = append(cars, c)
	}
	if err := rows.Err(); err != nil {
			return cars, err
	}
	return cars, nil
}

func GetCarByID(id int) (model.Car, error) {
	var car model.Car
	err := DB.QueryRow("SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars WHERE id = $1", id).Scan(
			&car.ID,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.Color,
			&car.BodyStyle,
			&car.EngineSize,
			&car.Weight,
			&car.BasePrice,
			&car.FuelCapacity,
			&car.Horsepower,
			&car.Torque,
			&car.Acceleration,
			&car.TopSpeed,
	)
	if err != nil {
			return car, err
	}
	return car, nil
}

func UpdateCarByID(id int, c model.Car) error {
	_, err := DB.Exec("UPDATE cars SET brand = $1, model = $2, year = $3, color = $4, body_style = $5, engine_size = $6, weight = $7, base_price = $8, fuel_capacity = $9, horsepower = $10, torque = $11, acceleration = $12, top_speed = $13 WHERE id = $14",
			c.Brand,
			c.Model,
			c.Year,
			c.Color,
			c.BodyStyle,
			c.EngineSize,
			c.Weight,
			c.BasePrice,
			c.FuelCapacity,
			c.Horsepower,
			c.Torque,
			c.Acceleration,
			c.TopSpeed,
			id,
	)
	if err != nil {
			return err
	}
	return nil
}

func DeleteCarByID(id int) error {
	_, err := DB.Exec("DELETE FROM cars WHERE id = $1", id)
	if err != nil {
			return err
	}
	return nil
}
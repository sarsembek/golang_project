package model

import (
	"database/sql"
	"errors"
)

type Car struct {
	ID           int     `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Color        string  `json:"color"`
	BodyStyle    string  `json:"body_style"`
	EngineSize   float64 `json:"engine_size"`
	Weight       float64 `json:"weight"`
	BasePrice    int     `json:"base_price"`
	FuelCapacity int     `json:"fuel_capacity"`
	Horsepower   int     `json:"horsepower"`
	Torque       int     `json:"torque"`
	Acceleration int     `json:"acceleration"`
	TopSpeed     int     `json:"top_speed"`
}

// CreateCar inserts a new car into the database
func CreateCar(db *sql.DB, c Car) error {
	_, err := db.Exec("INSERT INTO cars (brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
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

// GetAllCars retrieves all cars from the database
func GetAllCars(db *sql.DB) ([]Car, error) {
	var cars []Car
	rows, err := db.Query("SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Car
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
			return nil, err
		}
		cars = append(cars, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

// GetCarWithPagination retrieves cars with pagination, filtering, and sorting
func GetCarWithPagination(db *sql.DB, page, limit int, sortBy, filterBy string) ([]Car, error) {
	var cars []Car

	query := "SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars"

	if filterBy != "" {
		query += " WHERE brand LIKE '%" + filterBy + "%' OR model LIKE '%" + filterBy + "%'"
	}

	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	query += " LIMIT $1 OFFSET $2"
	offset := (page - 1) * limit

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Car
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
			return nil, err
		}
		cars = append(cars, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

// GetCarByID retrieves a car by ID from the database
func GetCarByID(db *sql.DB, id int) (*Car, error) {
	var car Car
	err := db.QueryRow("SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars WHERE id = $1", id).
		Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Color, &car.BodyStyle, &car.EngineSize, &car.Weight, &car.BasePrice, &car.FuelCapacity, &car.Horsepower, &car.Torque, &car.Acceleration, &car.TopSpeed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil, nil if no rows found
		}
		return nil, err
	}
	return &car, nil
}

// UpdateCarByID updates a car by ID in the database
func UpdateCarByID(db *sql.DB, id int, c Car) error {
	_, err := db.Exec("UPDATE cars SET brand = $1, model = $2, year = $3, color = $4, body_style = $5, engine_size = $6, weight = $7, base_price = $8, fuel_capacity = $9, horsepower = $10, torque = $11, acceleration = $12, top_speed = $13 WHERE id = $14",
		c.Brand, c.Model, c.Year, c.Color, c.BodyStyle, c.EngineSize, c.Weight, c.BasePrice, c.FuelCapacity, c.Horsepower, c.Torque, c.Acceleration, c.TopSpeed, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCarByID deletes a car by ID from the database
func DeleteCarByID(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

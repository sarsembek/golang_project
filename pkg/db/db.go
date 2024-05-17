package db

import (
	"database/sql"
	"errors"
	"log"

	"car_project/pkg/model"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err := sql.Open("postgres", "user=postgres password=postgres dbname=cars host=localhost port=5432 sslmode=disable")
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
		"file://migrations", // Replace with the path to your migration files
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

// CreateCar inserts a new car into the database
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

// GetCarWithPagination retrieves cars with pagination, filtering, and sorting
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
		return nil, err
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
func GetCarByID(id int) (*model.Car, error) {
	var car model.Car
	err := DB.QueryRow("SELECT id, brand, model, year, color, body_style, engine_size, weight, base_price, fuel_capacity, horsepower, torque, acceleration, top_speed FROM cars WHERE id = $1", id).
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
func UpdateCarByID(id int, c model.Car) error {
	_, err := DB.Exec("UPDATE cars SET brand = $1, model = $2, year = $3, color = $4, body_style = $5, engine_size = $6, weight = $7, base_price = $8, fuel_capacity = $9, horsepower = $10, torque = $11, acceleration = $12, top_speed = $13 WHERE id = $14",
		c.Brand, c.Model, c.Year, c.Color, c.BodyStyle, c.EngineSize, c.Weight, c.BasePrice, c.FuelCapacity, c.Horsepower, c.Torque, c.Acceleration, c.TopSpeed, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCarByID deletes a car by ID from the database
func DeleteCarByID(id int) error {
	_, err := DB.Exec("DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// CreateCarHistory inserts a new car history record into the database
func CreateCarHistory(carHistory model.CarHistory) error {
    _, err := DB.Exec("INSERT INTO car_history (car_id, date, type, details, service_type, service_cost, service_notes) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        carHistory.CarID, carHistory.Date, carHistory.Type, carHistory.Details, carHistory.ServiceType, carHistory.ServiceCost, carHistory.ServiceNotes)
    if err != nil {
        return err
    }
    return nil
}
// GetCarHistoryWithPagination retrieves car history with pagination, filtering, and sorting
func GetCarAllHistory(page int, limit int, sortBy, filterBy string) ([]model.CarHistory, error) {
    // Construct SQL query based on pagination, filtering, and sorting parameters
    query := "SELECT * FROM car_history"

    // Add filtering if specified
    if filterBy != "" {
        query += " WHERE " + filterBy
    }

    // Add sorting if specified
    if sortBy != "" {
        query += " ORDER BY " + sortBy
    }

    // Add pagination
    query += " LIMIT $1 OFFSET $2"

    // Execute query
    rows, err := DB.Query(query, limit, (page-1)*limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Parse rows into CarHistory objects
    var carHistory []model.CarHistory
    for rows.Next() {
        var ch model.CarHistory
        err := rows.Scan(&ch.ID, &ch.CarID, &ch.Date, &ch.Type, &ch.Details, &ch.ServiceType, &ch.ServiceCost, &ch.ServiceNotes)
        if err != nil {
            return nil, err
        }
        carHistory = append(carHistory, ch)
    }

    return carHistory, nil
}

// GetCarHistoryByID retrieves a car history record by ID from the database
func GetCarHistoryByID(id int) (model.CarHistory, error) {
    var carHistory model.CarHistory
    err := DB.QueryRow("SELECT id, car_id, date, type, details, service_type, service_cost, service_notes FROM car_history WHERE id = $1", id).
        Scan(&carHistory.ID, &carHistory.CarID, &carHistory.Date, &carHistory.Type, &carHistory.Details, &carHistory.ServiceType, &carHistory.ServiceCost, &carHistory.ServiceNotes)
    if err != nil {
        return model.CarHistory{}, err
    }
    return carHistory, nil
}

// UpdateCarHistory updates an existing car history record in the database
func UpdateCarHistory(carHistory model.CarHistory) error {
    _, err := DB.Exec("UPDATE car_history SET car_id = $1, date = $2, type = $3, details = $4, service_type = $5, service_cost = $6, service_notes = $7 WHERE id = $8",
        carHistory.CarID, carHistory.Date, carHistory.Type, carHistory.Details, carHistory.ServiceType, carHistory.ServiceCost, carHistory.ServiceNotes, carHistory.ID)
    if err != nil {
        return err
    }
    return nil
}

// DeleteCarHistory deletes a car history record by ID from the database
func DeleteCarHistory(id int) error {
    _, err := DB.Exec("DELETE FROM car_history WHERE id = $1", id)
    if err != nil {
        return err
    }
    return nil
}
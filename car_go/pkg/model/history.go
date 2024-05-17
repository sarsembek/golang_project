package model

import (
	"database/sql"
	"time"
)

type CarHistory struct {
    ID           int       `json:"id"`
    CarID        int       `json:"car_id"`
    Date         time.Time `json:"date"`
    Type         string    `json:"type"` // Type of event: "accident" or "service"
    Details      string    `json:"details"`
    ServiceType  string    `json:"service_type,omitempty"`  // Type of service (e.g., oil change, maintenance)
    ServiceCost  float64   `json:"service_cost,omitempty"`  // Cost of service
    ServiceNotes string    `json:"service_notes,omitempty"` // Additional notes for service
}

// CreateCarHistory inserts a new car history record into the database
func CreateCarHistory(db *sql.DB, carHistory CarHistory) error {
    _, err := db.Exec("INSERT INTO car_history (car_id, date, type, details, service_type, service_cost, service_notes) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        carHistory.CarID, carHistory.Date, carHistory.Type, carHistory.Details, carHistory.ServiceType, carHistory.ServiceCost, carHistory.ServiceNotes)
    if err != nil {
        return err
    }
    return nil
}
// GetCarHistoryWithPagination retrieves car history with pagination, filtering, and sorting
func GetCarAllHistory(db *sql.DB, page int, limit int, sortBy, filterBy string) ([]CarHistory, error) {
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
    rows, err := db.Query(query, limit, (page-1)*limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Parse rows into CarHistory objects
    var carHistory []CarHistory
    for rows.Next() {
        var ch CarHistory
        err := rows.Scan(&ch.ID, &ch.CarID, &ch.Date, &ch.Type, &ch.Details, &ch.ServiceType, &ch.ServiceCost, &ch.ServiceNotes)
        if err != nil {
            return nil, err
        }
        carHistory = append(carHistory, ch)
    }

    return carHistory, nil
}

// GetCarHistoryByID retrieves a car history record by ID from the database
func GetCarHistoryByID(db *sql.DB, id int) (CarHistory, error) {
    var carHistory CarHistory
    err := db.QueryRow("SELECT id, car_id, date, type, details, service_type, service_cost, service_notes FROM car_history WHERE id = $1", id).
        Scan(&carHistory.ID, &carHistory.CarID, &carHistory.Date, &carHistory.Type, &carHistory.Details, &carHistory.ServiceType, &carHistory.ServiceCost, &carHistory.ServiceNotes)
    if err != nil {
        return CarHistory{}, err
    }
    return carHistory, nil
}

// UpdateCarHistory updates an existing car history record in the database
func UpdateCarHistory(db *sql.DB, carHistory CarHistory) error {
    _, err := db.Exec("UPDATE car_history SET car_id = $1, date = $2, type = $3, details = $4, service_type = $5, service_cost = $6, service_notes = $7 WHERE id = $8",
        carHistory.CarID, carHistory.Date, carHistory.Type, carHistory.Details, carHistory.ServiceType, carHistory.ServiceCost, carHistory.ServiceNotes, carHistory.ID)
    if err != nil {
        return err
    }
    return nil
}

// DeleteCarHistory deletes a car history record by ID from the database
func DeleteCarHistory(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM car_history WHERE id = $1", id)
    if err != nil {
        return err
    }
    return nil
}
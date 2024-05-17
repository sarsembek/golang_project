package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"car_project/pkg/model"
)

// CreateCarHistory creates a new car history record
func CreateCarHistory(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var carHistory model.CarHistory
        err := json.NewDecoder(r.Body).Decode(&carHistory)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate input fields
        if carHistory.CarID <= 0 {
            http.Error(w, "Invalid CarID", http.StatusBadRequest)
            return
        }

        // Check if the car with the provided ID exists
        car, err := model.GetCarByID(db, carHistory.CarID)
        if err != nil {
            http.Error(w, "Failed to get car by ID", http.StatusInternalServerError)
            return
        }
        if car == nil {
            http.Error(w, "CarID does not exist", http.StatusBadRequest)
            return
        }
        if carHistory.Type != "accident" && carHistory.Type != "service" {
            http.Error(w, "Invalid Type", http.StatusBadRequest)
            return
        }

        // Additional input validation logic can be added here

        carHistory.Date = time.Now() // Set current time as the date
        err = model.CreateCarHistory(db, carHistory)
        if err != nil {
            http.Error(w, "Failed to create car history", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func GetAllCarHistory(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        page, _ := strconv.Atoi(query.Get("page"))
        limit, _ := strconv.Atoi(query.Get("limit"))
        sortBy := query.Get("sortBy")
        filterBy := query.Get("filterBy")

        if page == 0 {
            page = 1
        }
        if limit == 0 {
            limit = 10
        }

        carHistory, err := model.GetCarAllHistory(db, page, limit, sortBy, filterBy)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(carHistory)
    }
}

// GetCarHistoryByID retrieves a car history record by ID
func GetCarHistoryByID(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            http.Error(w, "Invalid car history ID", http.StatusBadRequest)
            return
        }

        carHistory, err := model.GetCarHistoryByID(db, id)
        if err != nil {
            http.Error(w, "Failed to get car history", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(carHistory)
    }
}

// UpdateCarHistory updates an existing car history record
func UpdateCarHistory(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var carHistory model.CarHistory
        err := json.NewDecoder(r.Body).Decode(&carHistory)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        err = model.UpdateCarHistory(db, carHistory)
        if err != nil {
            http.Error(w, "Failed to update car history", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

// DeleteCarHistory deletes a car history record by ID
func DeleteCarHistory(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            http.Error(w, "Invalid car history ID", http.StatusBadRequest)
            return
        }

        err = model.DeleteCarHistory(db, id)
        if err != nil {
            http.Error(w, "Failed to delete car history", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

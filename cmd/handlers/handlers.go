package handlers

import (
	"car_project/pkg/db"
	"car_project/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("YourSecretKey"), nil // Use your actual secret key
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Inject user ID into the context for downstream handlers to use
			type contextKey string

			const (
				userIDKey contextKey = "userID"
			)
			ctx := context.WithValue(r.Context(), userIDKey, claims["userID"])
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Store the user in the database (pseudo-code)
	err = db.CreateUser(user)
	if err != nil {
		http.Error(w, "Error storing user in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials model.User
	_ = json.NewDecoder(r.Body).Decode(&credentials)

	// Retrieve user from the database
	user, err := db.GetUserByUsername(credentials.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userID":   user.ID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("YourSecretKey"))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func Activate(w http.ResponseWriter, r *http.Request) {
	// Extract the expired token from the request header or request body
	expiredToken := r.Header.Get("Expired-Token")
	if expiredToken == "" {
		http.Error(w, "Expired token not provided", http.StatusBadRequest)
		return
	}

	// Parse the expired token
	token, err := jwt.Parse(expiredToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("YourSecretKey"), nil // Use your actual secret key
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract claims from the expired token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userID, ok := claims["userID"].(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Generate a new JWT token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := newToken.SignedString([]byte("YourSecretKey"))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the new JWT token
	w.Write([]byte(tokenString))
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
		var c model.Car
		_ = json.NewDecoder(r.Body).Decode(&c)
		err := db.CreateCar(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
}

func GetAllCars(w http.ResponseWriter, r *http.Request) {
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

		cars, err := db.GetCarWithPagination(page, limit, sortBy, filterBy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cars)
	}

func GetCar(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		car, err := db.GetCarByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(car)
	}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars((r))
		id, _ := strconv.Atoi(params["id"])
		var c model.Car
		_ = json.NewDecoder(r.Body).Decode(&c)
		err := db.UpdateCarByID(id, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		err := db.DeleteCarByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

func CreateCarHistory(w http.ResponseWriter, r *http.Request) {
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
		car, err := db.GetCarByID(carHistory.CarID)
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
		err = db.CreateCarHistory(carHistory)
		if err != nil {
			http.Error(w, "Failed to create car history", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

func GetAllCarHistory(w http.ResponseWriter, r *http.Request) {
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

		carHistory, err := db.GetCarAllHistory(page, limit, sortBy, filterBy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(carHistory)
	}

// GetCarHistoryByID retrieves a car history record by ID
func GetCarHistoryByID(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid car history ID", http.StatusBadRequest)
			return
		}

		carHistory, err := db.GetCarHistoryByID(id)
		if err != nil {
			http.Error(w, "Failed to get car history", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(carHistory)
	}


// UpdateCarHistory updates an existing car history record
func UpdateCarHistory(w http.ResponseWriter, r *http.Request) {
		var carHistory model.CarHistory
		err := json.NewDecoder(r.Body).Decode(&carHistory)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = db.UpdateCarHistory(carHistory)
		if err != nil {
			http.Error(w, "Failed to update car history", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}


// DeleteCarHistory deletes a car history record by ID
func DeleteCarHistory(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid car history ID", http.StatusBadRequest)
			return
		}

		err = db.DeleteCarHistory(id)
		if err != nil {
			http.Error(w, "Failed to delete car history", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
}

func CreateRating(w http.ResponseWriter, r *http.Request) {
    var rating model.Rating
    err := json.NewDecoder(r.Body).Decode(&rating)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate car_id
    exists, err := db.CarExists(rating.CarID)
    if err != nil {
        http.Error(w, "Error checking car ID", http.StatusInternalServerError)
        return
    }
    if !exists {
        http.Error(w, "Car ID does not exist", http.StatusBadRequest)
        return
    }

    // Validate user_id
    if rating.UserID <= 0 {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    // Insert the rating into the database
    // (Assuming you have a ratings table)
    insertQuery := "INSERT INTO ratings (car_id, stars, user_id, comment) VALUES ($1, $2, $3, $4)"
    _, err = db.DB.Exec(insertQuery, rating.CarID, rating.Stars, rating.UserID, rating.Comment)
    if err != nil {
        http.Error(w, "Error inserting rating into database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}


// GetRating retrieves the rating for a specific car
func GetRating(w http.ResponseWriter, r *http.Request) {
    // Parse the car_id from the query parameters
    carIDParam := r.URL.Query().Get("car_id")
    carID, err := strconv.Atoi(carIDParam)
    if err != nil {
        http.Error(w, "Invalid car_id", http.StatusBadRequest)
        return
    }

    // Validate car_id
    exists, err := db.CarExists(carID)
    if err != nil {
        http.Error(w, "Error checking car ID", http.StatusInternalServerError)
        return
    }
    if !exists {
        http.Error(w, "Car ID does not exist", http.StatusBadRequest)
        return
    }

    // Query the database for the rating
    query := "SELECT car_id, stars, user_id, comment FROM ratings WHERE car_id = $1"
    row := db.DB.QueryRow(query, carID)

    var rating model.Rating
    err = row.Scan(&rating.CarID, &rating.Stars, &rating.UserID, &rating.Comment)
    if err != nil {
        http.Error(w, "Rating not found", http.StatusNotFound)
        return
    }

    // Return the rating as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rating)
}


func UpdateRating(w http.ResponseWriter, r *http.Request) {
    carID := r.URL.Query().Get("car_id")
    userID := r.URL.Query().Get("user_id")

    var updatedRating model.Rating
    err := json.NewDecoder(r.Body).Decode(&updatedRating)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Convert carID and userID to integers
    carIDInt, err := strconv.Atoi(carID)
    if err != nil {
        http.Error(w, "Invalid car_id", http.StatusBadRequest)
        return
    }
    userIDInt, err := strconv.Atoi(userID)
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    // Update the rating in the database
    err = db.UpdateRating(carIDInt, userIDInt, updatedRating)
    if err != nil {
        http.Error(w, "Error updating rating in database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}


func DeleteRating(w http.ResponseWriter, r *http.Request) {
    carID := r.URL.Query().Get("car_id")
    userID := r.URL.Query().Get("user_id")

    // Convert carID and userID to integers
    carIDInt, err := strconv.Atoi(carID)
    if err != nil {
        http.Error(w, "Invalid car_id", http.StatusBadRequest)
        return
    }
    userIDInt, err := strconv.Atoi(userID)
    if err != nil {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    // Delete the rating from the database
    err = db.DeleteRating(carIDInt, userIDInt)
    if err != nil {
        http.Error(w, "Error deleting rating from database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

package handlers

import (
    "fmt"
	"encoding/json"
	"net/http"
	"strconv"
    "strings"
    "context"
    "time"

    "golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"car_project/pkg/model"
	"car_project/pkg/db"
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
			return []byte("YourSecretKey"), nil  // Use your actual secret key
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Inject user ID into the context for downstream handlers to use
			ctx := context.WithValue(r.Context(), "userID", claims["userID"])
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
        "userID": user.ID,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString([]byte("YourSecretKey"))
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

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
    params := mux.Vars(r)
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

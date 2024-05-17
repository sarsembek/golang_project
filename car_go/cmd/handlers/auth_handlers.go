package handlers

import (
	"car_project/pkg/db"
	"car_project/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
			return []byte("YourSecretKey"), nil  // Use your actual secret key
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
        return []byte("YourSecretKey"), nil  // Use your actual secret key
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
        "exp": time.Now().Add(time.Hour * 72).Unix(),
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

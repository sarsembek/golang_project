package token

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJWT(w http.ResponseWriter, r *http.Request) {
    // Authentication logic here

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": "exampleUser",  // Payload data
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString([]byte("YourSecretKey"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(tokenString))
}


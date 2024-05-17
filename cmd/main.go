package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "car_project/cmd/handlers"
    "car_project/pkg/db"
)

func main() {
    // Initialize the database
    db.InitDB()

    // Create a new router
    r := mux.NewRouter()
    api := r.PathPrefix("/api").Subrouter()
    api.Use(handlers.Authenticate)

    // Define routes
    api.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
    api.HandleFunc("/cars", handlers.GetAllCars).Methods("GET")
    api.HandleFunc("/cars/{id}", handlers.GetCar).Methods("GET")
    api.HandleFunc("/cars/{id}", handlers.UpdateCar).Methods("PUT")
    api.HandleFunc("/cars/{id}", handlers.DeleteCar).Methods("DELETE")

    api.HandleFunc("/carhistory", handlers.CreateCarHistory).Methods("POST") 
    api.HandleFunc("/carhistory", handlers.GetAllCarHistory).Methods("GET") 
    api.HandleFunc("/carhistory/{id}", handlers.GetCarHistoryByID).Methods("GET") 
    api.HandleFunc("/carhistory/{id}", handlers.UpdateCarHistory).Methods("PUT") 
    api.HandleFunc("/carhistory/{id}", handlers.DeleteCarHistory).Methods("DELETE") 

    r.HandleFunc("/user/register", handlers.RegisterUser).Methods("POST")
    r.HandleFunc("/user/login", handlers.LoginUser).Methods("POST")
    r.HandleFunc("user/activate", handlers.Activate).Methods("POST")

    // Start the server
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

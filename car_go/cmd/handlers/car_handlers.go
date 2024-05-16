package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"car_project/pkg/model"
)


func CreateCar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var c model.Car
        _ = json.NewDecoder(r.Body).Decode(&c)
        err := model.CreateCar(db, c)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAllCars(db *sql.DB) http.HandlerFunc {
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

        cars, err := model.GetCarWithPagination(db, page, limit, sortBy, filterBy)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(cars)
    }
}

func GetCar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        id, _ := strconv.Atoi(params["id"])
        car, err := model.GetCarByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(car)
    }
}

func UpdateCar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        id, _ := strconv.Atoi(params["id"])
        var c model.Car
        _ = json.NewDecoder(r.Body).Decode(&c)
        err := model.UpdateCarByID(db, id, c)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}

func DeleteCar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        id, _ := strconv.Atoi(params["id"])
        err := model.DeleteCarByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}
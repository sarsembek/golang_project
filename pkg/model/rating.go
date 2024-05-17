package model

type Rating struct {
  CarID   int      `json:"car_id"`
  Stars   int      `json:"stars"`
  UserID   int      `json:"user_id"`
  Comment string   `json:"comment"`
}
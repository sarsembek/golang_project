package model

type Car struct {
	ID           int     `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Color        string  `json:"color"`
	BodyStyle    string  `json:"body_style"`
	EngineSize   float64 `json:"engine_size"`
	Weight       float64 `json:"weight"`
	BasePrice    int     `json:"base_price"`
	FuelCapacity int     `json:"fuel_capacity"`
	Horsepower   int     `json:"horsepower"`
	Torque       int     `json:"torque"`
	Acceleration int     `json:"acceleration"`
	TopSpeed     int     `json:"top_speed"`
}

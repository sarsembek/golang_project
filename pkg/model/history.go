package model

import (
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

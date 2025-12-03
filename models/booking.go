package models

import "time"

type Booking struct {
	ID          string    `json:"id"`
	Movie       string    `json:"movie"`
	MovieNumber string    `json:"movie_number"`
	Seat        string    `json:"seat"`
	User        string    `json:"user"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Runtime     int       `json:"runtime"`
	MPAARating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Genre       string    `json:"genre,omitempty"`
	Image       string    `json:"image,omitempty"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

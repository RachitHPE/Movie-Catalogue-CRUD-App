package model

import "time"

type MovieInfo struct {
	ID        int64     `gorm:"primaryKey" json:"id"` // Unique integer ID for movies
	CreatedAt time.Time `json:"createdAt"`            // Timestamp for creation of a movie
	UpdatedAt time.Time `json:"updatedAt"`            // Timestamp for updation of a movie
	Title     string    `json:"title"`                // String title for movie
	Year      int32     `json:"year,omitempty"`       // Movie release year
	Genres    []string  `json:"genres,omitempty"`     // Slice of genres for the movie
}

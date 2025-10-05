package models

import (
	"strconv"
	"time"
)

type Ascent struct {
	Number      uint      `json:"number" bson:"number"`
	Date        string    `json:"date" bson:"date"`
	Country     string    `json:"country" bson:"country"`
	CountryCode string    `json:"countryCode" bson:"countryCode"`
	Area        string    `json:"area" bson:"area"`
	City        string    `json:"city" bson:"city"`
	Crag        string    `json:"crag" bson:"crag"`
	Sector      string    `json:"sector" bson:"sector"`
	Name        string    `json:"name" bson:"name"`
	Grade       string    `json:"grade" bson:"grade"`
	Style       string    `json:"style" bson:"style"`
	CreatedAt   time.Time `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"-" bson:"updatedAt,omitempty"`
}

type PublicAscent struct {
	Number      uint   `json:"number"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Area        string `json:"area"`
	City        string `json:"city"`
	Crag        string `json:"crag"`
	Sector      string `json:"sector"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	Style       string `json:"style"`
	Year        int    `json:"year"`
}

// serializeAscent converts an Ascent to PublicAscent
func SerializeAscent(ascent Ascent) PublicAscent {
	year := 0
	// Parse year from date string (assuming format like "DD/MM/YYYY")
	if len(ascent.Date) >= 10 {
		if y, err := strconv.Atoi(ascent.Date[6:10]); err == nil {
			year = y
		}
	}

	return PublicAscent{
		Number:      ascent.Number,
		Country:     ascent.Country,
		CountryCode: ascent.CountryCode,
		Area:        ascent.Area,
		City:        ascent.City,
		Crag:        ascent.Crag,
		Sector:      ascent.Sector,
		Name:        ascent.Name,
		Grade:       ascent.Grade,
		Style:       ascent.Style,
		Year:        year,
	}
}

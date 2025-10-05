package models

import "time"

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

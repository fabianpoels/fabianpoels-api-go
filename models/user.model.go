package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,alpha"`
	Password  string             `json:"-" bson:"password,omitempty" validate:"required,alpha"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,alpha"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedAt,omitempty"`
}

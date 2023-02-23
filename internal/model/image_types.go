package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	CatId    string    `bson:"catId,omitempty" json:"catId,omitempty"`
	ImageUrl string    `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
}

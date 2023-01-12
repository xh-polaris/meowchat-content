package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cat struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Age          string             `bson:"age,omitempty"`
	CommunityId  string             `bson:"communityId,omitempty"`
	Color        string             `bson:"color,omitempty"`
	Details      string             `bson:"details,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Sex          string             `bson:"sex,omitempty"`
	Status       int64              `bson:"status,omitempty"`
	Area         string             `bson:"area,omitempty"`
	IsSnipped    bool               `bson:"isSnipped,omitempty"`
	IsSterilized bool               `bson:"isSterilized,omitempty"`
	Avatars      []string           `bson:"avatars,omitempty"`
	UpdateAt     time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt     time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

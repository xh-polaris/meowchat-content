package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cat struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Age          string             `bson:"age"`
	CommunityId  string             `bson:"communityId"`
	Color        string             `bson:"color"`
	Details      string             `bson:"details"`
	Name         string             `bson:"name"`
	Sex          string             `bson:"sex"`
	Status       int64              `bson:"status"`
	Area         string             `bson:"area"`
	IsSnipped    bool               `bson:"isSnipped"`
	IsSterilized bool               `bson:"isSterilized"`
	IsDeleted    bool               `bson:"isDeleted"`
	Avatars      []string           `bson:"avatars"`
	DeleteAt     time.Time          `bson:"deleteAt,omitempty"`
	UpdateAt     time.Time          `bson:"updateAt" json:"updateAt"`
	CreateAt     time.Time          `bson:"createAt" json:"createAt"`
}

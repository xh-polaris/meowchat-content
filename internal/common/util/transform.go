package util

import (
	"encoding/json"
	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	"strconv"
	"time"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
)

func BoolToInt(a bool) int64 {
	if a {
		return 1
	}
	return 0
}
func IntToBool(a int64) bool {
	return a == 1
}

func TransformPbCat(Cat *model.Cat) *pb.Cat {
	id := strconv.FormatInt(Cat.Id, 10)
	var s []string
	err := json.Unmarshal([]byte(Cat.Avatars), &s)
	if err != nil {
		return nil
	}
	return &pb.Cat{
		Id:           id,
		CreateAt:     Cat.CreateAt.Unix(),
		Age:          Cat.Age,
		CommunityId:  Cat.CommunityId,
		Color:        Cat.Color,
		Details:      Cat.Details,
		Name:         Cat.Name,
		Popularity:   Cat.Popularity,
		Sex:          Cat.Sex,
		Status:       Cat.Status,
		Area:         Cat.Area,
		IsSnipped:    IntToBool(Cat.IsSnipped),
		IsSterilized: IntToBool(Cat.IsSterilized),
		Avatars:      s,
	}
}

func TransformModelCat(Cat *pb.Cat) (*model.Cat, error) {
	id, err := strconv.ParseInt(Cat.Id, 10, 64)
	if err != nil {
		if Cat.Id == "" {
			id = 0
		} else {
			return nil, errorx.ErrInvalidId
		}
	}
	str, err := json.Marshal(Cat.Avatars)
	if err != nil {
		return nil, err
	}
	return &model.Cat{
		Id:           id,
		CreateAt:     time.Unix(Cat.CreateAt, 0),
		Age:          Cat.Age,
		CommunityId:  Cat.CommunityId,
		Color:        Cat.Color,
		Details:      Cat.Details,
		Name:         Cat.Name,
		Popularity:   Cat.Popularity,
		Sex:          Cat.Sex,
		Status:       Cat.Status,
		Area:         Cat.Area,
		IsSnipped:    BoolToInt(Cat.IsSnipped),
		IsSterilized: BoolToInt(Cat.IsSterilized),
		Avatars:      string(str),
	}, nil
}

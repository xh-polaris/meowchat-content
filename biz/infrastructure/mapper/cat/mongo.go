package cat

import (
	"context"
	"time"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const prefixCatCacheKey = "cache:cat:"
const CollectionName = "cat"

var _ IMongoMapper = (*MongoMapper)(nil)

type (
	IMongoMapper interface {
		Insert(ctx context.Context, data *Cat) error
		FindOne(ctx context.Context, id string) (*Cat, error)
		Update(ctx context.Context, data *Cat) error
		Delete(ctx context.Context, id string) error
		FindManyByCommunityId(ctx context.Context, communityId string, skip int64, count int64) ([]*Cat, int64, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	Cat struct {
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
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Cat) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixCatCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Cat, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data Cat
	key := prefixCatCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{consts.ID: oid})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) Update(ctx context.Context, data *Cat) error {
	data.UpdateAt = time.Now()
	key := prefixCatCacheKey + data.ID.Hex()
	_, err := m.conn.UpdateOne(ctx, key, bson.M{consts.ID: data.ID}, bson.M{"$set": data})
	return err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixCatCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

func (m *MongoMapper) FindManyByCommunityId(ctx context.Context, communityId string, skip int64, count int64) ([]*Cat, int64, error) {
	data := make([]*Cat, 0, 20)
	err := m.conn.Find(ctx, &data, bson.M{consts.CommunityId: communityId}, &options.FindOptions{
		Skip:  &skip,
		Limit: &count,
		Sort:  bson.M{consts.ID: -1},
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := m.conn.CountDocuments(ctx, bson.M{consts.CommunityId: communityId})
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

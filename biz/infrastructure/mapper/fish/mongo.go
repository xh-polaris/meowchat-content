package fish

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
)

var _ IMongoMapper = (*MongoMapper)(nil)

const prefixFishCacheKey = "cache:fish:"
const CollectionName = "fish"

type (
	IMongoMapper interface {
		Insert(ctx context.Context, data *Fish) error
		FindOne(ctx context.Context, id string) (*Fish, error)
		Update(ctx context.Context, data *Fish) error
		Delete(ctx context.Context, id string) error
		Add(ctx context.Context, id string, add int64) error
		StartClient(ctx context.Context) (*mongo.Client, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	Fish struct {
		UserId   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
		FishNum  int64              `bson:"fishNum,omitempty" json:"fishNum,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Fish) error {
	data.CreateAt = time.Now()
	data.UpdateAt = time.Now()
	key := prefixFishCacheKey + data.UserId.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Fish, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data Fish
	key := prefixFishCacheKey + id
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

func (m *MongoMapper) Update(ctx context.Context, data *Fish) error {
	data.UpdateAt = time.Now()
	key := prefixFishCacheKey + data.UserId.Hex()
	_, err := m.conn.ReplaceOne(ctx, key, bson.M{consts.ID: data.UserId}, data)
	return err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixFishCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

func (m *MongoMapper) Add(ctx context.Context, id string, add int64) error {
	key := prefixFishCacheKey + id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	filter := bson.M{consts.ID: oid}
	update := bson.M{"$inc": bson.M{consts.FishNum: add}, "$set": bson.M{consts.UpdateAt: time.Now()}}
	_, err = m.conn.UpdateOne(ctx, key, filter, update)
	return err
}

func (m *MongoMapper) StartClient(ctx context.Context) (*mongo.Client, error) {
	client := m.conn.Database().Client()
	return client, nil
}

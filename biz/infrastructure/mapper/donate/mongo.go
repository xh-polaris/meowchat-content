package donate

import (
	"context"
	"time"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
)

var _ IMongoMapper = (*MongoMapper)(nil)

const prefixDonateCacheKey = "cache:donate:"
const CollectionName = "donate"

type (
	IMongoMapper interface {
		Insert(ctx context.Context, data *Donate) error
		FindOne(ctx context.Context, id string) (*Donate, error)
		Update(ctx context.Context, data *Donate) error
		Delete(ctx context.Context, id string) error
		ListDonateByPlan(ctx context.Context, planId string) ([]*Donate, error)
		ListDonateByUser(ctx context.Context, userId string) ([]*Donate, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	Donate struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
		UserId   string             `bson:"userId,omitempty" json:"userId,omitempty"`
		PlanId   string             `bson:"planId,omitempty" json:"planId,omitempty"`
		FishNum  int64              `bson:"fishNum,omitempty" json:"fishNum,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Donate) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixDonateCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Donate, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data Donate
	key := prefixDonateCacheKey + id
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

func (m *MongoMapper) Update(ctx context.Context, data *Donate) error {
	data.UpdateAt = time.Now()
	key := prefixDonateCacheKey + data.ID.Hex()
	_, err := m.conn.ReplaceOne(ctx, key, bson.M{consts.ID: data.ID}, data)
	return err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixDonateCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

func (m *MongoMapper) ListDonateByPlan(ctx context.Context, planId string) ([]*Donate, error) {
	var data []*Donate
	var err error

	opts := options.FindOptions{
		Sort: bson.M{consts.FishNum: -1},
	}
	filter := bson.M{consts.PlanId: planId}

	err = m.conn.Find(ctx, &data, filter, &opts)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (m *MongoMapper) ListDonateByUser(ctx context.Context, userId string) ([]*Donate, error) {
	var data []*Donate
	var err error

	opts := options.FindOptions{
		Sort: bson.M{consts.CreateAt: -1},
	}
	filter := bson.M{consts.UserId: userId}

	err = m.conn.Find(ctx, &data, filter, &opts)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

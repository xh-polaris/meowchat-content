package donate

import (
	"context"
	"github.com/samber/lo"
	"sync"
	"time"

	"github.com/xh-polaris/gopkg/pagination"
	"github.com/xh-polaris/gopkg/pagination/mongop"

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
		FindManyAndCountByUserId(ctx context.Context, id string, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Donate, int64, error)
		FindManyByUserId(ctx context.Context, id string, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Donate, error)
		CountByUserId(ctx context.Context, id string) (int64, error)
		FindManyAndCountByPlanId(ctx context.Context, id string, popts *pagination.PaginationOptions) ([]*Donate, int64, error)
		FindManyByPlanId(ctx context.Context, id string, popts *pagination.PaginationOptions) ([]*Donate, error)
		CountByPlanId(ctx context.Context, id string) (int64, error)
		FindOneById(ctx context.Context, planId string, userId string) (*Donate, error)
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

func (m *MongoMapper) FindManyByUserId(ctx context.Context, id string, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Donate, error) {
	p := mongop.NewMongoPaginator(pagination.NewRawStore(sorter), popts)

	filter := bson.M{consts.UserId: id}

	sort, err := p.MakeSortOptions(ctx, filter)
	if err != nil {
		return nil, err
	}

	var data []*Donate
	if err = m.conn.Find(ctx, &data, filter, &options.FindOptions{
		Sort:  sort,
		Limit: popts.Limit,
		Skip:  popts.Offset,
	}); err != nil {
		return nil, err
	}

	// 如果是反向查询，反转数据
	if *popts.Backward {
		lo.Reverse(data)
	}
	if len(data) > 0 {
		err = p.StoreCursor(ctx, data[0], data[len(data)-1])
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (m *MongoMapper) CountByUserId(ctx context.Context, id string) (int64, error) {
	filter := bson.M{consts.UserId: id}
	return m.conn.CountDocuments(ctx, filter)
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

func (m *MongoMapper) FindOneById(ctx context.Context, planId string, userId string) (*Donate, error) {
	var data Donate
	err := m.conn.FindOneNoCache(ctx, &data, bson.M{consts.PlanId: planId, consts.UserId: userId})
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

func (m *MongoMapper) FindManyAndCountByUserId(ctx context.Context, id string, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Donate, int64, error) {
	var data []*Donate
	var total int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer wg.Done()
		var err error
		data, err = m.FindManyByUserId(ctx, id, popts, sorter)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		total, err = m.CountByUserId(ctx, id)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		wg.Wait()
		defer close(c)
	}()
	if err := <-c; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (m *MongoMapper) FindManyAndCountByPlanId(ctx context.Context, id string, popts *pagination.PaginationOptions) ([]*Donate, int64, error) {
	var data []*Donate
	var total int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer wg.Done()
		var err error
		data, err = m.FindManyByPlanId(ctx, id, popts)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		total, err = m.CountByPlanId(ctx, id)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		wg.Wait()
		defer close(c)
	}()
	if err := <-c; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (m *MongoMapper) FindManyByPlanId(ctx context.Context, id string, popts *pagination.PaginationOptions) ([]*Donate, error) {

	filter := bson.M{consts.PlanId: id}
	sort := bson.D{
		bson.E{Key: consts.FishNum, Value: -1},
	}

	var data []*Donate
	if err := m.conn.Find(ctx, &data, filter, &options.FindOptions{
		Sort:  sort,
		Limit: popts.Limit,
		Skip:  popts.Offset,
	}); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *MongoMapper) CountByPlanId(ctx context.Context, id string) (int64, error) {
	filter := bson.M{consts.PlanId: id}
	return m.conn.CountDocuments(ctx, filter)
}

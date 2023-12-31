package moment

import (
	"context"
	"github.com/samber/lo"
	"sync"
	"time"

	"github.com/xh-polaris/gopkg/pagination"
	"github.com/xh-polaris/gopkg/pagination/mongop"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
)

const prefixMomentCacheKey = "cache:moment:"
const CollectionName = "moment"

type (
	IMongoMapper interface {
		Insert(ctx context.Context, data *Moment) error
		FindOne(ctx context.Context, id string) (*Moment, error)
		Update(ctx context.Context, data *Moment) error
		Delete(ctx context.Context, id string) error
		FindMany(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Moment, error)
		Count(ctx context.Context, filter *FilterOptions) (int64, error)
		FindManyAndCount(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Moment, int64, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	Moment struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		CatId       string             `bson:"catId,omitempty"`
		CommunityId string             `bson:"communityId,omitempty"`
		Photos      []string           `bson:"photos,omitempty"`
		Title       string             `bson:"title,omitempty"`
		Text        string             `bson:"text,omitempty"`
		UserId      string             `bson:"userId,omitempty"`
		UpdateAt    time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt    time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
		// 仅ES查询时使用
		Score_ float64 `bson:"_score,omitempty" json:"_score,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) FindMany(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Moment, error) {
	p := mongop.NewMongoPaginator(pagination.NewRawStore(sorter), popts)

	filter := makeMongoFilter(fopts)
	sort, err := p.MakeSortOptions(ctx, filter)
	if err != nil {
		return nil, err
	}

	var data []*Moment
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

func (m *MongoMapper) Count(ctx context.Context, filter *FilterOptions) (int64, error) {
	f := makeMongoFilter(filter)
	return m.conn.CountDocuments(ctx, f)
}

func (m *MongoMapper) FindManyAndCount(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Moment, int64, error) {
	var data []*Moment
	var total int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer wg.Done()
		var err error
		data, err = m.FindMany(ctx, fopts, popts, sorter)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		total, err = m.Count(ctx, fopts)
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

func (m *MongoMapper) Insert(ctx context.Context, data *Moment) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixMomentCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Moment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data Moment
	key := prefixMomentCacheKey + id
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

func (m *MongoMapper) Update(ctx context.Context, data *Moment) error {
	data.UpdateAt = time.Now()
	key := prefixMomentCacheKey + data.ID.Hex()
	_, err := m.conn.UpdateByID(ctx, key, data.ID, bson.M{"$set": data})
	return err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixMomentCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

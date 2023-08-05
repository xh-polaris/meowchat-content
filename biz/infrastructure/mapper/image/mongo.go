package image

import (
	"context"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"math"
	"time"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ IMongoMapper = (*MongoMapper)(nil)

const prefixImageCacheKey = "cache:image:"
const CollectionName = "image"

type (
	IMongoMapper interface {
		Insert(ctx context.Context, data *Image) error
		FindOne(ctx context.Context, id string) (*Image, error)
		Update(ctx context.Context, data *Image) error
		Delete(ctx context.Context, id string) error
		ListImage(ctx context.Context, catId string, lastId *string, limit, offset int64, backward bool) ([]*Image, error)
		InsertMany(ctx context.Context, image []*Image) error
		CountImage(ctx context.Context, catId string) (int64, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	Image struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
		CatId    string             `bson:"catId,omitempty" json:"catId,omitempty"`
		ImageUrl string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Image) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixImageCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Image, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data Image
	key := prefixImageCacheKey + id
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

func (m *MongoMapper) Update(ctx context.Context, data *Image) error {
	data.UpdateAt = time.Now()
	key := prefixImageCacheKey + data.ID.Hex()
	_, err := m.conn.ReplaceOne(ctx, key, bson.M{consts.ID: data.ID}, data)
	return err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixImageCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

func (m *MongoMapper) InsertMany(ctx context.Context, image []*Image) error {
	for i := 0; i < len(image); i++ {
		if image[i].ID.IsZero() {
			image[i].ID = primitive.NewObjectID()
			image[i].CreateAt = time.Now()
			image[i].UpdateAt = time.Now()
		}
	}
	data := make([]interface{}, len(image))
	for i := 0; i < len(image); i++ {
		data[i] = image[i]
	}
	_, err := m.conn.InsertMany(ctx, data)
	return err
}

func (m *MongoMapper) ListImage(ctx context.Context, catId string, lastId *string, limit, offset int64, backward bool) ([]*Image, error) {
	var data []*Image
	var oid primitive.ObjectID
	var err error

	// 构造lastId
	if lastId == nil {
		if backward {
			oid = primitive.NewObjectIDFromTimestamp(time.Unix(math.MinInt32, 0))
		} else {
			oid = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt32, 0))
		}
	} else {
		oid, err = primitive.ObjectIDFromHex(*lastId)
		if err != nil {
			return nil, consts.ErrInvalidObjectId
		}
	}

	// 构造请求，新的数据在前面，数值越大越新
	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		Sort:  bson.M{consts.ID: -1},
	}
	filter := bson.M{consts.CatId: catId}
	if backward {
		filter[consts.ID] = bson.M{"$gt": oid}
	} else {
		filter[consts.ID] = bson.M{"$lt": oid}
	}

	err = m.conn.Find(ctx, &data, filter, &opts)
	if err != nil {
		return nil, err
	} else if len(data) <= 0 {
		return data, nil
	} else if backward {
		return data, nil
	} else {
		return data, nil
	}
}

func (m *MongoMapper) CountImage(ctx context.Context, catId string) (int64, error) {
	total, err := m.conn.CountDocuments(ctx, bson.M{consts.CatId: catId})
	if err != nil {
		return 0, err
	}
	return total, nil
}

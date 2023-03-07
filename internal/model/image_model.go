package model

import (
	"context"
	"math"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ImageModel = (*customImageModel)(nil)

const ImageCollectionName = "image"

type (
	// ImageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImageModel.
	ImageModel interface {
		imageModel
		ListImage(ctx context.Context, catId string, lastId *string, limit, offset int64, backward bool) ([]*Image, error)
		InsertMany(ctx context.Context, image []*Image) error
		CountImage(ctx context.Context, catId string) (int64, error)
	}

	customImageModel struct {
		*defaultImageModel
	}
)

func (c *customImageModel) InsertMany(ctx context.Context, image []*Image) error {
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
	_, err := c.conn.InsertMany(ctx, data)
	return err
}

func (c *customImageModel) ListImage(ctx context.Context, catId string, lastId *string, limit, offset int64, backward bool) ([]*Image, error) {
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
			return nil, ErrInvalidObjectId
		}
	}

	// 构造请求，新的数据在前面，数值越大越新
	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		Sort:  bson.M{"_id": -1},
	}
	filter := bson.M{"catId": catId}
	if backward {
		filter["_id"] = bson.M{"$gt": oid}
	} else {
		filter["_id"] = bson.M{"$lt": oid}
	}

	err = c.conn.Find(ctx, &data, filter, &opts)
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

func (c *customImageModel) CountImage(ctx context.Context, catId string) (int64, error) {
	total, err := c.conn.CountDocuments(ctx, bson.M{"catId": catId})
	if err != nil {
		return 0, err
	}
	return total, nil
}

// NewImageModel returns a model for the mongo.
func NewImageModel(url, db string, c cache.CacheConf) ImageModel {
	conn := monc.MustNewModel(url, db, ImageCollectionName, c)
	return &customImageModel{
		defaultImageModel: newDefaultImageModel(conn),
	}
}

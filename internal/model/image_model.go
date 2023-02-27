package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mon"
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
		ListImageByCat(ctx context.Context, catId string, lastId primitive.ObjectID, limit int64) ([]*Image, error)
		InsertMany(ctx context.Context, image []*Image) error
	}

	customImageModel struct {
		*defaultImageModel
	}
)

func (c customImageModel) InsertMany(ctx context.Context, image []*Image) error {
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

func (c customImageModel) ListImageByCat(ctx context.Context, catId string, lastId primitive.ObjectID, limit int64) ([]*Image, error) {

	var data []*Image

	opts := options.FindOptions{
		Limit: &limit,
		Sort:  bson.D{bson.E{Key: "_id", Value: -1}},
	}
	filter := bson.M{"catId": catId, "_id": bson.M{"$lt": lastId}}
	err := c.conn.Find(ctx, &data, filter, &opts)
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewImageModel returns a model for the mongo.
func NewImageModel(url, db string, c cache.CacheConf) ImageModel {
	conn := monc.MustNewModel(url, db, ImageCollectionName, c)
	return &customImageModel{
		defaultImageModel: newDefaultImageModel(conn),
	}
}

package model

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CatCollectionName = "cat"

var _ CatModel = (*customCatModel)(nil)

type (
	// CatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCatModel.
	CatModel interface {
		catModel
		FindManyByCommunityId(ctx context.Context, communityId string, skip int64, count int64) ([]*Cat, int64, error)
		Search(ctx context.Context, communityId, keyword string, skip, count int64) ([]*Cat, int64, error)
	}

	customCatModel struct {
		*defaultCatModel
		es        *elasticsearch.Client
		indexName string
	}
)

// NewCatModel returns a model for the mongo.
func NewCatModel(url, db string, c cache.CacheConf, es config.ElasticsearchConf) CatModel {
	conn := monc.MustNewModel(url, db, CatCollectionName, c)
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  es.Username,
		Password:  es.Password,
		Addresses: es.Addresses,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return &customCatModel{
		defaultCatModel: newDefaultCatModel(conn),
		es:              esClient,
		indexName:       fmt.Sprintf("%s.%s-alias", db, CatCollectionName),
	}
}

func (m *customCatModel) FindManyByCommunityId(ctx context.Context, communityId string, skip int64, count int64) ([]*Cat, int64, error) {
	data := make([]*Cat, 0, 20)
	err := m.conn.Find(ctx, &data, bson.M{"communityId": communityId}, &options.FindOptions{
		Skip:  &skip,
		Limit: &count,
		Sort:  bson.M{"createAt": -1},
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := m.conn.CountDocuments(ctx, bson.M{"communityId": communityId})
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (m *customCatModel) Search(ctx context.Context, communityId, keyword string, skip, count int64) ([]*Cat, int64, error) {
	search := m.es.Search
	query := map[string]any{
		"from": skip,
		"size": count,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"term": map[string]any{
							"communityId": communityId,
						},
					},
					map[string]any{
						"multi_match": map[string]any{
							"query":  keyword,
							"fields": []string{"details", "name", "area", "color"},
						},
					},
				},
			},
		},
		"sort": map[string]any{
			"_score": map[string]any{
				"order": "desc",
			},
			"createAt": map[string]any{
				"order": "desc",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}
	res, err := search(
		search.WithIndex(m.indexName),
		search.WithContext(ctx),
		search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, err
		} else {
			logx.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]any
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	hits := r["hits"].(map[string]any)["hits"].([]any)
	total := int64(r["hits"].(map[string]any)["total"].(map[string]any)["value"].(float64))
	cats := make([]*Cat, 0, 10)
	for i := range hits {
		hit := hits[i].(map[string]any)
		cat := &Cat{}
		source := hit["_source"].(map[string]any)
		if source["createAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["createAt"].(string)); err != nil {
			return nil, 0, err
		}
		if source["updateAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["updateAt"].(string)); err != nil {
			return nil, 0, err
		}
		hit["_source"] = source
		err := mapstructure.Decode(hit["_source"], cat)
		if err != nil {
			return nil, 0, err
		}
		oid := hit["_id"].(string)
		id, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return nil, 0, err
		}
		cat.ID = id
		cats = append(cats, cat)
	}
	return cats, total, nil
}

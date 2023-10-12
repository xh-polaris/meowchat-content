package plan

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
)

type FilterOptions struct {
	OnlyUserId *string
	OnlyCatId  *string
}

type MongoFilter struct {
	m bson.M
	*FilterOptions
}

func makeMongoFilter(options *FilterOptions) bson.M {
	return (&MongoFilter{
		m:             bson.M{},
		FilterOptions: options,
	}).toBson()
}

func (f *MongoFilter) toBson() bson.M {
	f.CheckOnlyUserId()
	f.CheckOnlyCatId()
	return f.m
}

func (f *MongoFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.m[consts.InitiatorIds] = *f.OnlyUserId
	}
}

func (f *MongoFilter) CheckOnlyCatId() {
	if f.OnlyCatId != nil {
		f.m[consts.CatId] = *f.OnlyCatId
	}
}

type EsFilter struct {
	q []types.Query
	*FilterOptions
}

func makeEsFilter(opts *FilterOptions) []types.Query {
	return (&EsFilter{
		q:             make([]types.Query, 0),
		FilterOptions: opts,
	}).toQuery()
}

func (f *EsFilter) toQuery() []types.Query {
	f.checkOnlyUserId()
	f.checkOnlyCatId()
	return f.q
}

func (f *EsFilter) checkOnlyUserId() {
	if f.OnlyUserId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.InitiatorIds: {Value: *f.OnlyUserId},
			},
		})
	}
}

func (f *EsFilter) checkOnlyCatId() {
	if f.OnlyCatId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.CatId: {Value: *f.OnlyCatId},
			},
		})
	}
}

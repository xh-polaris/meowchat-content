package moment

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"go.mongodb.org/mongo-driver/bson"
)

type FilterOptions struct {
	OnlyUserId       *string
	OnlyCommunityId  *string
	OnlyCommunityIds []string
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
	f.CheckOnlyCommunityId()
	f.CheckOnlyCommunityIds()
	return f.m
}

func (f *MongoFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.m[consts.UserId] = *f.OnlyUserId
	}
}

func (f *MongoFilter) CheckOnlyCommunityId() {
	if f.OnlyCommunityId != nil {
		f.m[consts.CommunityId] = *f.OnlyCommunityId
	}
}

func (f *MongoFilter) CheckOnlyCommunityIds() {
	if f.OnlyCommunityIds != nil {
		f.m[consts.CommunityId] = f.OnlyCommunityIds
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
	f.checkOnlyCommunityId()
	f.checkOnlyCommunityIds()
	return f.q
}

func (f *EsFilter) checkOnlyUserId() {
	if f.OnlyUserId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.UserId: {Value: *f.OnlyUserId},
			},
		})
	}
}

func (f *EsFilter) checkOnlyCommunityId() {
	if f.OnlyCommunityId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.CommunityId: {Value: *f.OnlyCommunityId},
			},
		})
	}
}

func (f *EsFilter) checkOnlyCommunityIds() {
	if f.OnlyCommunityIds != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.CommunityId: {Value: f.OnlyCommunityIds},
			},
		})
	}
}

package plan

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
)

type FilterOptions struct {
	OnlyUserId      *string
	OnlyCommunityId *string
	OnlyCatId       *string
	IncludeGlobal   *bool
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
	f.CheckOnlyCommunityId()
	return f.m
}

func (f *MongoFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.m[consts.InitiatorId] = *f.OnlyUserId
	}
}

func (f *MongoFilter) CheckOnlyCatId() {
	if f.OnlyCatId != nil {
		f.m[consts.CatId] = *f.OnlyCatId
	}
}

func (f *MongoFilter) CheckOnlyCommunityId() {
	if f.IncludeGlobal == nil {
		if f.OnlyCommunityId != nil {
			f.m[consts.CommunityId] = *f.OnlyCommunityId
		}
	} else if *f.IncludeGlobal == false {
		if f.OnlyCommunityId != nil {
			f.m[consts.CommunityId] = *f.OnlyCommunityId
		}
	} else {
		if f.OnlyCommunityId != nil {
			f.m["$or"] = bson.A{bson.M{consts.CommunityId: bson.M{"$exists": false}}, bson.M{consts.CommunityId: *f.OnlyCommunityId}}
		} else {
			f.m[consts.CommunityId] = bson.M{"$exists": false}
		}
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
	f.checkOnlyCommunityId()
	return f.q
}

func (f *EsFilter) checkOnlyUserId() {
	if f.OnlyUserId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				consts.InitiatorId: {Value: *f.OnlyUserId},
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

func (f *EsFilter) checkOnlyCommunityId() {
	if f.IncludeGlobal == nil {
		if f.OnlyCommunityId != nil {
			f.q = append(f.q, types.Query{
				Term: map[string]types.TermQuery{
					consts.CommunityId: {Value: *f.OnlyCommunityId},
				},
			})
		}
	} else if *f.IncludeGlobal == false {
		if f.OnlyCommunityId != nil {
			f.q = append(f.q, types.Query{
				Term: map[string]types.TermQuery{
					consts.CommunityId: {Value: *f.OnlyCommunityId},
				},
			})
		}
	} else {
		if f.OnlyCommunityId != nil {
			BoolQuery := make([]types.Query, 0)
			BoolQuery = append(BoolQuery, types.Query{
				Bool: &types.BoolQuery{
					MustNot: []types.Query{
						types.Query{
							Exists: &types.ExistsQuery{
								Field: consts.CommunityId,
							},
						},
					},
				},
			})
			BoolQuery = append(BoolQuery, types.Query{
				Term: map[string]types.TermQuery{
					consts.CommunityId: {Value: *f.OnlyCommunityId},
				},
			})
			f.q = append(f.q, types.Query{
				Bool: &types.BoolQuery{
					Should: BoolQuery,
				},
			})
		} else {
			BoolQuery := make([]types.Query, 0)
			BoolQuery = append(BoolQuery, types.Query{
				Bool: &types.BoolQuery{
					MustNot: []types.Query{
						types.Query{
							Exists: &types.ExistsQuery{
								Field: consts.CommunityId,
							},
						},
					},
				},
			})
			f.q = append(f.q, types.Query{
				Bool: &types.BoolQuery{
					Should: BoolQuery,
				},
			})
		}
	}
}

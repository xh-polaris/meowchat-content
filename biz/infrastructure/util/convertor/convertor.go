package convertor

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xh-polaris/gopkg/pagination"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/moment"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/plan"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/post"
)

func ConvertMomentSlice(data []*moment.Moment) []*content.Moment {
	res := make([]*content.Moment, len(data))
	for i, d := range data {
		m := &content.Moment{
			Id:          d.ID.Hex(),
			CreateAt:    d.CreateAt.Unix(),
			Photos:      d.Photos,
			Title:       d.Title,
			Text:        d.Text,
			UserId:      d.UserId,
			CommunityId: d.CommunityId,
			CatId:       &d.CatId,
		}
		res[i] = m
	}
	return res
}

func ConvertMoment(data *moment.Moment) *content.Moment {
	return &content.Moment{
		Id:          data.ID.Hex(),
		CreateAt:    data.CreateAt.Unix(),
		Photos:      data.Photos,
		Title:       data.Title,
		Text:        data.Text,
		UserId:      data.UserId,
		CommunityId: data.CommunityId,
		CatId:       &data.CatId,
	}
}

func ConvertPlan(data *plan.Plan) *content.Plan {
	return &content.Plan{
		Id:          data.ID.Hex(),
		CreateAt:    data.CreateAt.Unix(),
		PlanType:    data.PlanType,
		EndTime:     data.EndTime.Unix(),
		Description: data.Description,
		StartTime:   data.StartTime.Unix(),
		InitiatorId: data.InitiatorId,
		CatId:       data.CatId,
		CommunityId: data.CommunityId,
		ImageUrls:   data.ImageUrls,
		Name:        data.Name,
		Instruction: data.Instruction,
		CoverUrl:    data.CoverUrl,
		Summary:     data.Summary,
		MaxFish:     data.MaxFish,
		NowFish:     data.NowFish,
		PlanState:   data.PlanState,
	}
}

func ConvertPost(in *post.Post) *content.Post {
	return &content.Post{
		Id:         in.ID.Hex(),
		CreateAt:   in.CreateAt.Unix(),
		UpdateAt:   in.UpdateAt.Unix(),
		Title:      in.Title,
		Text:       in.Text,
		CoverUrl:   in.CoverUrl,
		Tags:       in.Tags,
		UserId:     in.UserId,
		IsOfficial: in.Flags.GetFlag(post.OfficialFlag),
	}
}

func ConvertMomentAllFieldsSearchQuery(in *content.SearchOptions_AllFieldsKey) []types.Query {
	return []types.Query{{
		MultiMatch: &types.MultiMatchQuery{
			Query:  in.AllFieldsKey,
			Fields: []string{consts.Title + "^3", consts.Text},
		}},
	}
}

func ConvertMomentMultiFieldsSearchQuery(in *content.SearchOptions_MultiFieldsKey) []types.Query {
	var q []types.Query
	if in.MultiFieldsKey.Title != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Title: {
					Query: *in.MultiFieldsKey.Title + "^3",
				},
			},
		})
	}
	if in.MultiFieldsKey.Text != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Text: {
					Query: *in.MultiFieldsKey.Text,
				},
			},
		})
	}
	return q
}

func ConvertPlanAllFieldsSearchQuery(in *content.SearchOptions_AllFieldsKey) []types.Query {
	return []types.Query{{
		MultiMatch: &types.MultiMatchQuery{
			Query:  in.AllFieldsKey,
			Fields: []string{consts.Name + "^3", consts.Description},
		}},
	}
}

func ConvertPlanMultiFieldsSearchQuery(in *content.SearchOptions_MultiFieldsKey) []types.Query {
	var q []types.Query
	if in.MultiFieldsKey.Title != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Name: {
					Query: *in.MultiFieldsKey.Title + "^3",
				},
			},
		})
	}
	if in.MultiFieldsKey.Text != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Description: {
					Query: *in.MultiFieldsKey.Text,
				},
			},
		})
	}
	return q
}

func ConvertPostAllFieldsSearchQuery(in *content.SearchOptions_AllFieldsKey) []types.Query {
	return []types.Query{{
		MultiMatch: &types.MultiMatchQuery{
			Query:  in.AllFieldsKey,
			Fields: []string{consts.Title + "^3", consts.Text, consts.Tags},
		}},
	}
}

func ConvertPostMultiFieldsSearchQuery(in *content.SearchOptions_MultiFieldsKey) []types.Query {
	var q []types.Query
	if in.MultiFieldsKey.Title != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Title: {
					Query: *in.MultiFieldsKey.Title + "^3",
				},
			},
		})
	}
	if in.MultiFieldsKey.Text != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Text: {
					Query: *in.MultiFieldsKey.Text,
				},
			},
		})
	}
	if in.MultiFieldsKey.Tag != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				consts.Tags: {
					Query: *in.MultiFieldsKey.Tag,
				},
			},
		})
	}
	return q
}

func ParseMomentFilter(opts *content.MomentFilterOptions) (filter *moment.FilterOptions) {
	if opts == nil {
		filter = &moment.FilterOptions{}
	} else {
		filter = &moment.FilterOptions{
			OnlyUserId:       opts.OnlyUserId,
			OnlyCommunityId:  opts.OnlyCommunityId,
			OnlyCommunityIds: opts.OnlyCommunityIds,
		}
	}
	return
}

func ParsePlanFilter(opts *content.PlanFilterOptions) (filter *plan.FilterOptions) {
	if opts == nil {
		filter = &plan.FilterOptions{}
	} else {
		filter = &plan.FilterOptions{
			OnlyUserId:      opts.OnlyUserId,
			OnlyCommunityId: opts.OnlyCommunityId,
			OnlyCatId:       opts.OnlyCatId,
			IncludeGlobal:   opts.IncludeGlobal,
		}
	}
	return
}

func ParsePostFilter(fopts *content.PostFilterOptions) *post.FilterOptions {
	if fopts != nil {
		return &post.FilterOptions{
			OnlyUserId:   fopts.OnlyUserId,
			OnlyOfficial: fopts.OnlyOfficial,
		}
	}
	return &post.FilterOptions{}
}

func ParsePagination(opts *basic.PaginationOptions) (p *pagination.PaginationOptions) {
	if opts == nil {
		p = &pagination.PaginationOptions{}
	} else {
		p = &pagination.PaginationOptions{
			Limit:     opts.Limit,
			Offset:    opts.Offset,
			Backward:  opts.Backward,
			LastToken: opts.LastToken,
		}
	}
	return
}

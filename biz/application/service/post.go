package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/gopkg/pagination/esp"
	"github.com/xh-polaris/gopkg/pagination/mongop"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/post"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/convertor"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *content.CreatePostReq) (*content.CreatePostResp, error)
	RetrievePost(ctx context.Context, req *content.RetrievePostReq) (*content.RetrievePostResp, error)
	UpdatePost(ctx context.Context, req *content.UpdatePostReq) (*content.UpdatePostResp, error)
	DeletePost(ctx context.Context, req *content.DeletePostReq) (*content.DeletePostResp, error)
	ListPost(ctx context.Context, req *content.ListPostReq) (*content.ListPostResp, error)
	CountPost(ctx context.Context, req *content.CountPostReq) (*content.CountPostResp, error)
	SetOfficial(ctx context.Context, req *content.SetOfficialReq) (*content.SetOfficialResp, error)
}

type PostService struct {
	PostMongoMapper post.IMongoMapper
	PostEsMapper    post.IEsMapper
}

var PostSet = wire.NewSet(
	wire.Struct(new(PostService), "*"),
	wire.Bind(new(IPostService), new(*PostService)),
)

func (s *PostService) CreatePost(ctx context.Context, req *content.CreatePostReq) (*content.CreatePostResp, error) {
	p := &post.Post{
		Title:    req.Title,
		Text:     req.Text,
		CoverUrl: req.CoverUrl,
		Tags:     req.Tags,
		UserId:   req.UserId,
	}
	err := s.PostMongoMapper.Insert(ctx, p)
	if err != nil {
		return nil, err
	}
	return &content.CreatePostResp{PostId: p.ID.Hex()}, nil
}

func (s *PostService) RetrievePost(ctx context.Context, req *content.RetrievePostReq) (*content.RetrievePostResp, error) {
	data, err := s.PostMongoMapper.FindOne(ctx, req.PostId)
	switch err {
	case nil:
	case consts.ErrNotFound:
		return nil, consts.ErrNoSuchPost
	default:
		return nil, err
	}
	return &content.RetrievePostResp{Post: convertor.ConvertPost(data)}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *content.UpdatePostReq) (*content.UpdatePostResp, error) {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	err = s.PostMongoMapper.Update(ctx, &post.Post{
		ID:       oid,
		Title:    req.Title,
		Text:     req.Text,
		CoverUrl: req.CoverUrl,
		Tags:     req.Tags,
	})
	if err != nil {
		return nil, err
	}

	return &content.UpdatePostResp{}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *content.DeletePostReq) (*content.DeletePostResp, error) {
	err := s.PostMongoMapper.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &content.DeletePostResp{}, nil
}

func (s *PostService) ListPost(ctx context.Context, req *content.ListPostReq) (*content.ListPostResp, error) {
	resp := new(content.ListPostResp)
	var posts []*post.Post
	var total int64
	var err error

	filter := convertor.ParsePostFilter(req.FilterOptions)

	p := convertor.ParsePagination(req.PaginationOptions)

	if req.SearchOptions == nil {
		posts, total, err = s.PostMongoMapper.FindManyAndCount(ctx, filter, p, mongop.IdCursorType)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			posts, total, err = s.PostEsMapper.Search(ctx, convertor.ConvertPostAllFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
		case *content.SearchOptions_MultiFieldsKey:
			posts, total, err = s.PostEsMapper.Search(ctx, convertor.ConvertPostMultiFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Posts = make([]*content.Post, 0, len(posts))
	for _, post_ := range posts {
		resp.Posts = append(resp.Posts, convertor.ConvertPost(post_))
	}
	return resp, nil
}

func (s *PostService) CountPost(ctx context.Context, req *content.CountPostReq) (*content.CountPostResp, error) {
	var total int64
	var err error

	filter := convertor.ParsePostFilter(req.FilterOptions)

	if req.SearchOptions == nil {
		total, err = s.PostMongoMapper.Count(ctx, filter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			total, err = s.PostEsMapper.CountWithQuery(ctx, convertor.ConvertPostAllFieldsSearchQuery(o), filter)
		case *content.SearchOptions_MultiFieldsKey:
			total, err = s.PostEsMapper.CountWithQuery(ctx, convertor.ConvertPostMultiFieldsSearchQuery(o), filter)
		}
		if err != nil {
			return nil, err
		}
	}

	return &content.CountPostResp{Total: total}, nil
}

func (s *PostService) SetOfficial(ctx context.Context, req *content.SetOfficialReq) (*content.SetOfficialResp, error) {
	err := s.PostMongoMapper.UpdateFlags(ctx, req.PostId, map[post.Flag]bool{
		post.OfficialFlag: !req.IsRemove,
	})
	if err != nil {
		return nil, err
	}
	return &content.SetOfficialResp{}, nil
}

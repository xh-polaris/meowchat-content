package service

import (
	"context"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	catmapper "github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/cat"

	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICatService interface {
	SearchCat(ctx context.Context, req *content.SearchCatReq) (res *content.SearchCatResp, err error)
	ListCat(ctx context.Context, req *content.ListCatReq) (res *content.ListCatResp, err error)
	RetrieveCat(ctx context.Context, req *content.RetrieveCatReq) (res *content.RetrieveCatResp, err error)
	CreateCat(ctx context.Context, req *content.CreateCatReq) (res *content.CreateCatResp, err error)
	UpdateCat(ctx context.Context, req *content.UpdateCatReq) (res *content.UpdateCatResp, err error)
	DeleteCat(ctx context.Context, req *content.DeleteCatReq) (res *content.DeleteCatResp, err error)
}

type CatService struct {
	CatMongoMapper catmapper.IMongoMapper
	CatEsMapper    catmapper.IEsMapper
}

var CatSet = wire.NewSet(
	wire.Struct(new(CatService), "*"),
	wire.Bind(new(ICatService), new(*CatService)),
)

func (s *CatService) SearchCat(ctx context.Context, req *content.SearchCatReq) (res *content.SearchCatResp, err error) {
	data, total, err := s.CatEsMapper.Search(ctx, req.CommunityId, req.Keyword, int(req.Skip), int(req.Count))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var cats []*content.Cat
	for _, val := range data {
		cat := &content.Cat{}
		err = copier.Copy(cat, val)
		if err != nil {
			return nil, err
		}
		cat.Id = val.ID.Hex()
		cat.CreateAt = val.CreateAt.Unix()
		cats = append(cats, cat)
	}
	return &content.SearchCatResp{Cats: cats, Total: total}, nil
}

func (s *CatService) ListCat(ctx context.Context, req *content.ListCatReq) (res *content.ListCatResp, err error) {
	data, total, err := s.CatMongoMapper.FindManyByCommunityId(ctx, req.CommunityId, req.Skip, req.Count)
	if err != nil {
		return nil, err
	}
	var cats []*content.Cat
	for _, val := range data {
		cat := &content.Cat{}
		err = copier.Copy(cat, val)
		if err != nil {
			return nil, err
		}
		cat.Id = val.ID.Hex()
		cat.CreateAt = val.CreateAt.Unix()
		cats = append(cats, cat)
	}
	return &content.ListCatResp{Cats: cats, Total: total}, nil
}

func (s *CatService) RetrieveCat(ctx context.Context, req *content.RetrieveCatReq) (res *content.RetrieveCatResp, err error) {
	data, err := s.CatMongoMapper.FindOne(ctx, req.CatId)
	switch err {
	case nil:
	case consts.ErrNotFound:
		return nil, consts.ErrNoSuchCat
	default:
		return nil, err
	}
	cat := &content.Cat{}
	err = copier.Copy(cat, data)
	if err != nil {
		return nil, err
	}
	cat.Id = data.ID.Hex()
	cat.CreateAt = data.CreateAt.Unix()
	return &content.RetrieveCatResp{Cat: cat}, nil
}

func (s *CatService) CreateCat(ctx context.Context, req *content.CreateCatReq) (res *content.CreateCatResp, err error) {
	cat := &catmapper.Cat{}
	err = copier.Copy(cat, req.Cat)
	if err != nil {
		return nil, err
	}
	err = s.CatMongoMapper.Insert(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &content.CreateCatResp{CatId: cat.ID.Hex()}, nil
}

func (s *CatService) UpdateCat(ctx context.Context, req *content.UpdateCatReq) (res *content.UpdateCatResp, err error) {
	cat := &catmapper.Cat{}
	err = copier.Copy(cat, req.Cat)
	if err != nil {
		return nil, err
	}
	cat.ID, err = primitive.ObjectIDFromHex(req.Cat.Id)
	if err != nil {
		return nil, consts.ErrInvalidId
	}
	err = s.CatMongoMapper.Update(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &content.UpdateCatResp{}, nil
}

func (s *CatService) DeleteCat(ctx context.Context, req *content.DeleteCatReq) (res *content.DeleteCatResp, err error) {
	err = s.CatMongoMapper.Delete(ctx, req.CatId)
	if err != nil {
		return nil, err
	}
	return &content.DeleteCatResp{}, nil
}

package service

import (
	"context"

	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/data/db"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/mapper"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/collection"

	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatService interface {
	SearchCat(ctx context.Context, req *collection.SearchCatReq) (res *collection.SearchCatResp, err error)
	ListCat(ctx context.Context, req *collection.ListCatReq) (res *collection.ListCatResp, err error)
	RetrieveCat(ctx context.Context, req *collection.RetrieveCatReq) (res *collection.RetrieveCatResp, err error)
	CreateCat(ctx context.Context, req *collection.CreateCatReq) (res *collection.CreateCatResp, err error)
	UpdateCat(ctx context.Context, req *collection.UpdateCatReq) (res *collection.UpdateCatResp, err error)
	DeleteCat(ctx context.Context, req *collection.DeleteCatReq) (res *collection.DeleteCatResp, err error)
}

type CatServiceImpl struct {
	CatModel mapper.CatModel
}

var CatSet = wire.NewSet(
	wire.Struct(new(CatServiceImpl), "*"),
	wire.Bind(new(CatService), new(*CatServiceImpl)),
)

func (s *CatServiceImpl) SearchCat(ctx context.Context, req *collection.SearchCatReq) (res *collection.SearchCatResp, err error) {
	data, total, err := s.CatModel.Search(ctx, req.CommunityId, req.Keyword, req.Skip, req.Count)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var cats []*collection.Cat
	for _, val := range data {
		cat := &collection.Cat{}
		err = copier.Copy(cat, val)
		if err != nil {
			return nil, err
		}
		cat.Id = val.ID.Hex()
		cat.CreateAt = val.CreateAt.Unix()
		cats = append(cats, cat)
	}
	return &collection.SearchCatResp{Cats: cats, Total: total}, nil
}

func (s *CatServiceImpl) ListCat(ctx context.Context, req *collection.ListCatReq) (res *collection.ListCatResp, err error) {
	data, total, err := s.CatModel.FindManyByCommunityId(ctx, req.CommunityId, req.Skip, req.Count)
	if err != nil {
		return nil, err
	}
	var cats []*collection.Cat
	for _, val := range data {
		cat := &collection.Cat{}
		err = copier.Copy(cat, val)
		if err != nil {
			return nil, err
		}
		cat.Id = val.ID.Hex()
		cat.CreateAt = val.CreateAt.Unix()
		cats = append(cats, cat)
	}
	return &collection.ListCatResp{Cats: cats, Total: total}, nil
}

func (s *CatServiceImpl) RetrieveCat(ctx context.Context, req *collection.RetrieveCatReq) (res *collection.RetrieveCatResp, err error) {
	data, err := s.CatModel.FindOne(ctx, req.CatId)
	switch err {
	case nil:
	case mapper.ErrNotFound:
		return nil, consts.ErrNoSuchCat
	default:
		return nil, err
	}
	cat := &collection.Cat{}
	err = copier.Copy(cat, data)
	if err != nil {
		return nil, err
	}
	cat.Id = data.ID.Hex()
	cat.CreateAt = data.CreateAt.Unix()
	return &collection.RetrieveCatResp{Cat: cat}, nil
}

func (s *CatServiceImpl) CreateCat(ctx context.Context, req *collection.CreateCatReq) (res *collection.CreateCatResp, err error) {
	cat := &db.Cat{}
	err = copier.Copy(cat, req.Cat)
	if err != nil {
		return nil, err
	}
	err = s.CatModel.Insert(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &collection.CreateCatResp{CatId: cat.ID.Hex()}, nil
}

func (s *CatServiceImpl) UpdateCat(ctx context.Context, req *collection.UpdateCatReq) (res *collection.UpdateCatResp, err error) {
	cat := &db.Cat{}
	err = copier.Copy(cat, req.Cat)
	if err != nil {
		return nil, err
	}
	cat.ID, err = primitive.ObjectIDFromHex(req.Cat.Id)
	if err != nil {
		return nil, consts.ErrInvalidId
	}
	err = s.CatModel.Update(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &collection.UpdateCatResp{}, nil
}

func (s *CatServiceImpl) DeleteCat(ctx context.Context, req *collection.DeleteCatReq) (res *collection.DeleteCatResp, err error) {
	err = s.CatModel.Delete(ctx, req.CatId)
	if err != nil {
		return nil, err
	}
	return &collection.DeleteCatResp{}, nil
}

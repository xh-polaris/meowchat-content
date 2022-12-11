package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	_ CatModel = (*customCatModel)(nil)
)

type (
	// CatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCatModel.
	CatModel interface {
		catModel
		DeleteSoftly(ctx context.Context, id int64) error
		FindOneValid(ctx context.Context, id int64) (*Cat, error)
		FindManyValidByCommunityIdValid(ctx context.Context, CommunityId string, skip int64, count int64) ([]*Cat, error)
	}

	customCatModel struct {
		*defaultCatModel
	}
)

// NewCatModel returns a model for the database table.
func NewCatModel(conn sqlx.SqlConn, c cache.CacheConf) CatModel {
	return &customCatModel{
		defaultCatModel: newCatModel(conn, c),
	}
}

func (m *defaultCatModel) DeleteSoftly(ctx context.Context, id int64) error {
	meowchatCollectionRpcCatIdKey := fmt.Sprintf("%s%v", cacheMeowchatCollectionRpcCatIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `is_delete` = true, `delete_at` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, time.Now(), id)
	}, meowchatCollectionRpcCatIdKey)
	return err
}

func (m *defaultCatModel) FindOneValid(ctx context.Context, id int64) (*Cat, error) {
	meowchatCollectionRpcCatIdKey := fmt.Sprintf("%s%v", cacheMeowchatCollectionRpcCatIdPrefix, id)
	var resp Cat
	err := m.QueryRowCtx(ctx, &resp, meowchatCollectionRpcCatIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and `is_delete` = false limit 1", catRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, errorx.ErrNoSuchCat
	default:
		return nil, err
	}
}

func (m *defaultCatModel) FindManyValidByCommunityIdValid(ctx context.Context, CommunityId string, skip int64, count int64) ([]*Cat, error) {
	var resp []*Cat
	query := fmt.Sprintf("select %s from %s where `community_id` = ? and `is_delete` = 0 limit ?,?", catRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, CommunityId, skip, count)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

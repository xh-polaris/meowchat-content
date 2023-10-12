package consts

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"google.golang.org/grpc/status"
)

var (
	ErrNoSuchCat             = status.Error(10101, "no such cat")
	ErrInvalidId             = status.Error(10102, "invalid id")
	ErrNoSuchPost            = status.Error(10301, "no such post")
	ErrPaginatorTokenExpired = status.Error(10303, "paginator token has been expired")
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)

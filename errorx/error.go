package errorx

import "google.golang.org/grpc/status"

var (
	ErrNoSuchCat = status.Error(10101, "no such cat")
	ErrInvalidId = status.Error(10102, "invalid id")
)

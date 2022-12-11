package errorx

import "google.golang.org/grpc/status"

var (
	ErrNoSuchCat = status.Error(10001, "no such cat")
)

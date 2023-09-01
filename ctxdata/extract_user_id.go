package ctxdata

import (
	"context"
	"errors"
)

// CtxKeyJwtUserId get uid from ctx
type CtxKeyJwtUserId struct{}

func UserIdFromCxt(ctx context.Context) (string, error) {
	uid, ok := ctx.Value(CtxKeyJwtUserId{}).(string)
	if !ok {
		return "", errors.New("get uid from ctx failed")
	}
	return uid, nil
}

package ctxdata

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// CtxKeyJwtUserId get uid from ctx
type CtxKeyJwtUserId struct{}

func GetUidFromCtx(ctx context.Context) string {
	uid, ok := ctx.Value(CtxKeyJwtUserId{}).(string)
	if !ok {
		logx.WithContext(ctx).Errorf("get uid from ctx failed")
		return ""
	}
	return uid
}

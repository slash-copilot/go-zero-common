package ctxdata

import (
	"context"

	xerrors "github.com/slash-copilot/go-zero-common/errors"
	xhttp "github.com/slash-copilot/go-zero-common/http"
)

// CtxKeyJwtUserId get uid from ctx
type CtxKeyJwtUserId struct{}

func UserIdFromCxt(ctx context.Context) (string, error) {
	uid, ok := ctx.Value(CtxKeyJwtUserId{}).(string)
	if !ok {
		return "", &xerrors.CodeMsg{
			Code: xhttp.BusinessContextCodeError,
			Msg:  "user_id not found in ctx",
		}
	}
	return uid, nil
}

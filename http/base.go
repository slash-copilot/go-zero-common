package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/slash-copilot/go-zero-common/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// BaseResponse is the base response struct.
type BaseResponse[T any] struct {
	// Code represents the business code, not the http status code.
	Code string `json:"code"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	Msg string `json:"msg"`
	// Data represents the business data.
	Data T `json:"data,omitempty"`
}

// JsonBaseResponse writes v into w with http.StatusOK.
func JsonBaseResponse(w http.ResponseWriter, v any) {
	httpx.OkJson(w, wrapBaseResponse(v))
}

// JsonBaseResponseCtx writes v into w with http.StatusOK.
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	httpx.OkJsonCtx(ctx, w, wrapBaseResponse(v))
}

// JsonErrorResponse writes err into w.
func JsonErrorResponse(w http.ResponseWriter, code int, v any) {
	httpx.WriteJson(w, code, wrapBaseResponse(v))
}

// JsonErrorResponseCtx writes err into w.
func JsonErrorResponseCtx(ctx context.Context, w http.ResponseWriter, code int, v any) {
	httpx.WriteJsonCtx(ctx, w, code, wrapBaseResponse(v))
}
	
func wrapBaseResponse(v any) BaseResponse[any] {
	var resp BaseResponse[any]
	switch data := v.(type) {
	case *errors.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case errors.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case *status.Status:
		resp.Code = fmt.Sprintf("rpc:%d", data.Code())
		resp.Msg = data.Message()
	case error:
		resp.Code = BusinessDefaultCodeError
		resp.Msg = data.Error()
	default:
		resp.Code = BusinessMsgOk
		resp.Msg = BusinessMsgOk
		resp.Data = v
	}

	return resp
}

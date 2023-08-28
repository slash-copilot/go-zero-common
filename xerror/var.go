package xerror

import (
	"github.com/zeromicro/x/errors"
)

func NewErrCodeMsg(errCode int, errMsg string) *errors.CodeMsg {
	return &errors.CodeMsg{Code: errCode, Msg: errMsg}
}
func NewErrCode(errCode int) *errors.CodeMsg {
	return &errors.CodeMsg{Code: errCode, Msg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *errors.CodeMsg {
	return &errors.CodeMsg{Code: SERVER_COMMON_ERROR, Msg: errMsg}
}

func NewDBErrMsg(errMsg string) *errors.CodeMsg {
	return &errors.CodeMsg{Code: DB_ERROR, Msg: errMsg}
}

package xerror

// 成功返回
const OK int = 200

// err 400
const SERVER_ERROR_400 int = 400

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const SERVER_COMMON_ERROR int = 100001
const REUQEST_PARAM_ERROR int = 100002
const TOKEN_EXPIRE_ERROR int = 100003
const TOKEN_GENERATE_ERROR int = 100004

const DB_ERROR int = 100005
const DB_UPDATE_AFFECTED_ZERO_ERROR int = 100006

package main

// 定义枚举类型 ErrorCodes
type ErrorCodes int

const (
	Success        ErrorCodes = 0
	Error_Json     ErrorCodes = 1001 //Json解析错误
	RPCFailed      ErrorCodes = 1002 //RPC请求错误
	VarifyExpired  ErrorCodes = 1003 //验证码过期
	VarifyCodeErr  ErrorCodes = 1004 //验证码错误
	UserExist      ErrorCodes = 1005 //用户已经存在
	PasswdErr      ErrorCodes = 1006 //密码错误
	EmailNotMatch  ErrorCodes = 1007 //邮箱不匹配
	PasswdUpFailed ErrorCodes = 1008 //更新密码失败
	PasswdInvalid  ErrorCodes = 1009 //密码更新失败
	TokenInvalid   ErrorCodes = 1010 //Token失效
	UidInvalid     ErrorCodes = 1011 //uid无效
)

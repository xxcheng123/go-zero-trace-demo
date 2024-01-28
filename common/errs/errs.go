package errs

import "google.golang.org/grpc/codes"

/**
*	@Author: xxcheng
*	@Email developer@xxcheng.cn
*	@Date: 2024/1/28
 */

type Code uint32

const (
	NotFound         Code = 404
	Boom             Code = 666
	ParamParseFailed Code = 9998

	Finally Code = 9999
)

const (
	MySQLInternalError Code = 1001 + iota
	MySQLRegisterError
)

const (
	UserNotExist Code = 2000 + iota
	PasswordIncorrect
	UserExpired
)

var code2Message = map[Code]string{
	NotFound:         "未找到",
	Finally:          "内部错误",
	ParamParseFailed: "参数解析失败",
	Boom:             "这是一个错误，小心咯",

	MySQLInternalError: "数据库错误",
	MySQLRegisterError: "注册失败",

	UserNotExist:      "用户不存在",
	PasswordIncorrect: "密码错误",
	UserExpired:       "用户过期",
}

func (c Code) Error() string {
	if m, ok := code2Message[c]; ok {
		return m
	} else {
		return code2Message[Finally]
	}
}
func (c Code) Code() codes.Code {
	return codes.Code(c)
}

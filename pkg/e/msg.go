package e

var MsgFlags = map[int]string{
	SUCCESS:                                   "ok",
	ERROR:                                     "fail",
	INVALID_PARAMS:                            "请求参数错误",
	ERROR_NOT_EXIST_USERNAME:                  "用户名不存在",
	ERROR_CHECK_USERNAME:                      "检查用户名失败",
	ERROR_QUERY_USER_BY_USERNAME_AND_PASSWORD: "查询用户失败",
	ERROR_PASSWORD_NOT_CORRECT:                "密码不正确",
	ERROR_EXIST_USERNAME:                      "用户名已存在",
	ERROR_EXIST_PHONE:                         "手机号码已存在",
	ERROR_CHECK_PHONE:                         "检查手机号码失败",
	ERROR_EXIST_EMAIL:                         "邮箱地址已存在",
	ERROR_CHECK_EMAIL:                         "检查邮箱地址失败",
	ERROR_SAVE_USER:                           "保存用户失败",
	ERROR_QUERY_QUESTION:                      "查询问题失败",
	ERROR_NOT_EXIST_QUESTION:                  "该用户未设置找回密码问题",
	ERROR_CHECK_QUESTION_ANSWER:               "检查找回密码的问题和答案失败",
	ERROR_QUESTION_ANSWER_NOT_CORRECT:         "找回密码的问题和答案不匹配",
	ERROR_USER_NOT_LOGIN:                      "用户未登录",
	ERROR_QUERY_USER:                          "查询用户失败",
	ERROR_USER_NOT_EXIST:                      "用户不存在",
	ERROR_RESET_PASSWORD:                      "重置密码失败",
	ERROR_UPDATE_USER:                         "更新用户失败",
	ERROR_FORGET_RESET_PASSWORD_TOKEN:         "token不正确或已失效",
	ERROR_NOT_ADMIN:                           "当前登录用户不是管理员",
	ERROR_SAVE_CATEGORY:                       "保存分类失败",
	ERROR_CHECK_CATEGORY_BY_PARENT_ID:         "检查分类失败",
	ERROR_NOT_EXIST_CATEGORY:                  "分类不存在",
	ERROR_GET_CATEGORY:                        "查询分类失败",
	ERROR_UPDATE_CATEGORY:                     "更新分类失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

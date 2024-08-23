package constant

const (
	LOGIN_EMAIL_TIMES             = "LOGIN:EMAIL:TIMES:%s:%s"
	SEND_EMAIL_CODE               = "EMAIL:SEND:CODE:%s:%s"    // 修改密码发送的code
	SEND_EMAIL_COUNT_4_MODIFY_PWD = "EMAIL:SEND:CODE:COUNT:%s" // 修改密码发送的邮箱code的次数
	CHANGE_PWD_COUNT              = "EMAIL:CHANGE:PWD:COUNT:%s"
	REGISTER_LIMIT_RATE           = "REGISTER:LIMIT:RATE"
	SEND_EMAIL_CODE_LIMIT_RATE    = "SEND:EMAIL:CODE:LIMIT:RATE"
	CAPTCHA_IDENTIFYING           = "CAPTCHA:%s:LOGIN:I:%s"
	CAPTCHA_LOGIN_MODE            = "LOGIN"
	CAPTCHA_SEND_EMAIL_CODE_MODE  = "SEND_EMAIL_CODE:%s"
	CAPTCHA_CHANGE_PWD_MODE       = "CHANGE_PWD"
	CAPTCHA_REG_MODE              = "REGISTER"

	REDIS_COLLAB_USER  = "collab:users"
	REDIS_COLLAB_TOKEN = "collab:token"

	REDIS_COLLAB_UPDATES = "collab/%s/%s:updates" // roomid 类型

	ARTICLE_ID_SECRET = "article"
)

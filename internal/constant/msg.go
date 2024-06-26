package constant

const (
	SYS_ERR                 = "系统错误"
	PARAM_ERR               = "参数错误"
	TOKEN_NIL               = "认证信息为空"
	TOKEN_GEN_ERR           = "用户令牌生成失败"
	REFRESH_TOKEN_GEN_ERR   = "用户刷新令牌生成失败"
	TOKEN_EXPIRE_FORMAT_ERR = "过期时间格式化错误"
	TOKEN_EXPIRE            = "身份已过期,请重新登录"
	TOKEN_CHECK_ERR         = "令牌无法认证"
	TOKEN_MALFORMED_ERR     = "令牌格式错误"
	TOKEN_NOT_VALID_YET     = "令牌没有生效"

	USER_NOT_EXISTS                 = "用户不存在"
	CAPTCHA_COUNT_ERR               = "验证码次数已使用完成，需要等待恢复时间"
	CAPTCHA_KEY_INTERVAL_ERR        = "验证码需要等待间隔时间%s秒"
	CAPTCHA_NEED                    = "请验证"
	LOGIN_TIMES_ERR                 = "登录次数太频繁，1小时后再登录"
	RSA_LOAD_PRIVATE_KEY_ERR        = "无法加载私钥" //failed to parse PEM block containing the private key
	RSA_PARSE_PRIVATE_ERR           = "无法解析私钥"
	RSA_LOAD_PUBLIC_KEY_ERR         = "无法加载公钥"
	RSA_PARSE_PUBLIC_ERR            = "无法解析公钥"
	LOGIN_ACCOUNT_ERR               = "账户密码错误"
	LOGIN_EMAIL_NO_EXIST            = "邮箱不存在"
	EMAIL_ALREADY_EXISTS_ERR        = "邮箱已存在"
	LOGIN_PASSWORD_ERR              = "密码错误"
	LOGIN_PASSWORD_DECRYPT_ERR      = "密码解密失败"
	LOGIN_TIME_EXPIRE               = "认证已经过期，无法登录"
	SEND_EMAIL_CODE_ERR             = "获取邮箱验证码失败"
	SEND_EMAIL_CODE_NOT_EQUAL_ERR   = "邮箱验证码校验失败"
	SEND_EMAIL_POOL_ERR             = "通过连接池发送邮箱失败"
	SEND_EMAIL_CODE_EXPIR_ERR       = "验证码已发送，有效期%s分钟，请查看"
	SEND_EMAIL_CODE_LIMIT_COUNT_ERR = "发送验证码过于频繁"
	IS_PASSWORD_ERR                 = "密码长度只支持6到20位，需最少一个大写字母和特殊字符~!@#$%^&*()_+中的一个"
	IS_PASSWORD_CONFIRM_ERR         = "确认密码错误"
	TWO_NEW_PASSWORD_NOT_EQUAL_ERR  = "两密码不一致"
	CHANGE_PWD_ERR                  = "更新密码失败"
	CHANGE_PWD_TITLE                = "修改密码"
	CHANGE_PWD_CONTENT              = "您正在修改密码，验证码是: %s , 转发可能导致账号被盗。"
	REGISTER_LIMIT_RATE_ERR         = "注册次数过于频繁"
	CAPTCHA_DO_ERR                  = "请重新验证"

	ARTICLE_LIKE_ERR          = "点赞失败"
	ARTICLE_UNLIKE_ERR        = "取消点赞失败"
	ARTICLE_UPDATE_ERR        = "更新文章失败"
	ARTICLE_SAVE_ERR          = "保存文章失败"
	ARTICLE_AUTHOR_ERR        = "文章作者缺失"
	ARTICLE_CHECK_LIKE_ERR    = "检查是否点赞失败"
	ARTICLE_DRAFT_HISTORY_ERR = "文章保存到历史记录失败"

	ARTICLE_DARFT_PARAM_ERR = "获取文章历史编辑参数错误"

	GET_CURRENT_PATH_ERR = "获取当前路径失败"
	IMAGE_UPLOAD_ERR     = "上传图片失败"
	UPLOAD_DIR_ERR       = "获取上传路径失败"
	UPLOAD_EXCEED_ERR    = "一天内只能上传50次"
	FILE_PERMISSION_ERR  = "没有权限操作文件"
	CREATE_DIR_ERR       = "创建文件夹失败"
	OPEN_FILE_ERR        = "打开文件失败"
	FILE_NOT_EXIST_ERR   = "文件不存在"
	DECODE_IMG_ERR       = "解析图片失败"

	COLLAB_INVITE_USER_ERR  = "邀请用户失败"
	COLLAB_INVITE_TTL_ERR   = "邀请失效时间出错"
	COLLAB_INVITE_USER_NONE = "邀请用户为空"
	COLLAB_KICKOUT_USER_ERR = "删除协作用户失败"
	COLLAB_EXIT_ERR         = "退出协作失败"
	COLLAB_LIST_ERR         = "获取协作列表失败"
	COLLAB_TOKEN_ERR        = "生成地址失败"
)

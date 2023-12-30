package slideCaptcha

const (
	// CodeKeyPrefix 缓存key前缀
	CodeKeyPrefix = "RUNNING:CAPTCHA:%s"

	// BlockPuzzleCaptcha 滑动验证码服务标识
	BlockPuzzleCaptcha = "blockPuzzle"

	// ClickWordCaptcha 点击验证码服务标识
	ClickWordCaptcha = "clickWord"

	RedisCacheKey = "redis"
)

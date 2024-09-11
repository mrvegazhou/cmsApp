package timex

import (
	"github.com/spf13/cast"
	"time"
)

func IsValidTimestamp(timestamp int64) bool {
	if timestamp < 0 {
		return false
	}
	if len(cast.ToString(timestamp)) != 10 {
		return false
	}

	// 将时间戳转换为 time.Time 对象
	tm := time.Unix(timestamp, 0)

	// 检查转换后的时间是否在合理范围内
	// 例如，我们可以检查它是否在未来或过去的某个不合理的时间点
	if tm.Before(time.Unix(0, 0)) { // 检查是否在 Unix 纪元之前
		return false
	}
	if tm.After(time.Now().Add(24 * time.Hour * 365 * 100)) { // 检查是否在未来 100 年内
		return false
	}

	return true
}

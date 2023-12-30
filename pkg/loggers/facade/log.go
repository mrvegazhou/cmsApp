package facade

import (
	"cmsApp/configs"
	"context"
)

type Log interface {
	Info(context.Context, string, map[string]string)
	Error(context.Context, string, map[string]string)
}

func NewLogger(path string) (logger Log) {

	var logType string = configs.App.Base.LogMedia

	switch logType {

	case "redis":
		logger = newRedisLog(path)

	default:
		logger = newZaplog(path)
	}

	return
}

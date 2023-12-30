package facade

import (
	"cmsApp/pkg/loggers/newer"
	"context"
)

type RedisLog struct {
	logger *newer.RedisLogger
}

func newRedisLog(path string) *RedisLog {

	return &RedisLog{
		logger: newer.NewRedisLogger(path),
	}
}

func (rlog RedisLog) Info(ctx context.Context, msg string, info map[string]string) {
	value := ctx.Value("requestId")
	if value != nil {
		info["request_id"] = value.(string)
	}

	rlog.logger.Info(msg, info)
}

func (rlog RedisLog) Error(ctx context.Context, msg string, info map[string]string) {
	value := ctx.Value("requestId")
	if value != nil {
		info["request_id"] = value.(string)
	}

	rlog.logger.Error(msg, info)
}

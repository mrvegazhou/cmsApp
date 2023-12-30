package loggers

import (
	"context"

	"cmsApp/pkg/loggers/facade"
)

/*
* 通用info日志
 */
func LogInfo(ctx context.Context, path string, msg string, info map[string]string) {
	log := facade.NewLogger(path)
	log.Info(ctx, msg, info)
}

/*
* 通用error日志
 */
func LogError(ctx context.Context, path string, msg string, info map[string]string) {
	log := facade.NewLogger(path)
	log.Error(ctx, msg, info)
}

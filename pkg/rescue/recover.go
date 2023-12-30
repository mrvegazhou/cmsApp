package rescue

import (
	"cmsApp/pkg/loggers/facade"
	"context"
	"runtime/debug"
	"time"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		facade.NewLogger("Recover").Error(context.Background(), "[Recovery from panic]", map[string]string{
			"time":  time.Now().String(),
			"error": p.(string),
			"stack": string(debug.Stack()),
		})
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		facade.NewLogger("Recover").Error(ctx.(context.Context), "[Recovery from panic]", map[string]string{
			"time":  time.Now().String(),
			"error": p.(string),
			"stack": string(debug.Stack()),
		})
	}
}

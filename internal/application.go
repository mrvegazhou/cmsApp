package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cmsApp/configs"
	"cmsApp/internal/cron"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Route *gin.Engine
}

func (app Application) Run() {

	srv := &http.Server{
		Addr:    configs.App.Base.Host + ":" + configs.App.Base.Port,
		Handler: app.Route,
	}
	fmt.Print(configs.App.Base.Host + ":" + configs.App.Base.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit

	log.Println("Cron Close ...")
	openCron, cronCloseCtx := cron.GraceClose()
	if openCron {
		<-cronCloseCtx.Done()
	}

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}

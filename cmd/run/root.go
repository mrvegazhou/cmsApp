package run

import (
	"cmsApp/pkg/postgresqlx"
	"cmsApp/pkg/redisClient"
	"fmt"
	"io/ioutil"
	"log"

	"cmsApp/configs"
	"cmsApp/internal"
	"cmsApp/internal/cron"
	"cmsApp/internal/router"
	"cmsApp/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run app",
	Run:   runFunction,
}

var (
	configPath string
	crontab    string
	mode       string
)

func init() {
	CmdRun.Flags().StringVarP(&configPath, "config path", "c", "", "config path")
	CmdRun.Flags().StringVarP(&mode, "mode", "m", "debug", "debug or release")
	CmdRun.Flags().StringVarP(&crontab, "cron", "t", "open", "scheduled task control open or close")
}

func runFunction(cmd *cobra.Command, args []string) {
	var err error
	fmt.Println(111111)
	showLogo()

	//判断是否编译线上版本
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	//定时任务
	if crontab == "open" {
		cron.Init()
	}

	err = configs.Init(configPath)
	if err != nil {
		log.Fatalf("start fail:[Config Init] %s", err.Error())
	}

	err = web.Init()
	if err != nil {
		log.Fatalf("start fail:[Web Init] %s", err.Error())
	}

	err = redisClient.Init()
	if err != nil {
		log.Fatalf("start fail:[Redis Init] %s", err.Error())
	}

	err = postgresqlx.Init()
	if err != nil {
		log.Fatalf("start fail:[Mysql Init] %s", err.Error())
	}

	showPanel()

	r, err := router.Init()
	if err != nil {
		log.Fatalf("start fail:[Route Init] %s", err.Error())
	}

	app := internal.Application{}
	r.SetEngine(&app)
	app.Run()
}

func showLogo() {
	fmt.Println("   _____ _                   _           _       ")
	fmt.Println("  / ____(_)         /\\      | |         (_)      ")
	fmt.Println(" | |  __ _ _ __    /  \\   __| |_ __ ___  _ _ __  ")
	fmt.Println(" | | |_ | | '_ \\  / /\\ \\ / _` | '_ ` _ \\| | '_ \\ ")
	fmt.Println(" | |__| | | | | |/ _____\\ (_| | | | | | | | | | |")
	fmt.Println("  \\_____|_|_| |_/_/    \\_\\__,_|_| |_| |_|_|_| |_| \n")
}

func showPanel() {
	fmt.Println("App running at:")
	fmt.Printf("- Http Address:   %c[%d;%d;%dm%s%c[0m \n", 0x1B, 0, 40, 34, "http://"+configs.App.Base.Host+":"+configs.App.Base.Port, 0x1B)
	fmt.Println("")
}

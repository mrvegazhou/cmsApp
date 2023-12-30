package main

import (
	"cmsApp/configs"
	"cmsApp/internal"
	"cmsApp/internal/router"
	"cmsApp/pkg/postgresqlx"
	"cmsApp/pkg/redisClient"
	"cmsApp/web"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	release bool = true
)

// @title GinAdmin Api
// @version 1.0
// @description GinAdmin 示例项目

// @contact.name gphper
// @contact.url https://github.com/gphper/ginadmin

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:20011
// @basepath /api
func main() {
	// 设置时区
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("set location fail: %s", err.Error())
		os.Exit(1)
	}
	time.Local = local

	//var rootCmd = &cobra.Command{
	//	Use: "cmsApp",
	//}
	//rootCmd.AddCommand(run.CmdRun)
	//rootCmd.Execute()

	err = configs.Init("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp")
	if err != nil {
		fmt.Println(err, "---config err-----")
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

	r, err := router.Init()
	if err != nil {
		log.Fatalf("start fail:[Route Init] %s", err.Error())
	}

	app := internal.Application{}

	r.SetEngine(&app)
	app.Run()
}

package cron

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"cmsApp/configs"
	"cmsApp/pkg/redisClient"
	"cmsApp/pkg/utils/filesystem"
	gstrings "cmsApp/pkg/utils/strings"
)

var ctx = context.Background()

/**
* 创建目录
**/
func WriteLog() {

	var wg sync.WaitGroup

	date := time.Now().AddDate(0, 0, -1).Local().Format("20060102")

	pattern := "logs:" + date + ":*"
	keys, _ := redisClient.GetRedisClient().Keys(ctx, pattern).Result()

	for _, key := range keys {
		path := strings.ReplaceAll(key, ":", "/")

		file, err := filesystem.OpenFile(gstrings.JoinStr(configs.RootPath, "/", path, ".log"))
		if err == nil {
			wg.Add(1)
			go writeFile(key, file, &wg)
		}
	}

	wg.Wait()
}

/**
* 内容写入到文件中
**/
func writeFile(key string, file io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	var start int64 = 0
	var end int64 = 2
	for {
		logs, _ := redisClient.GetRedisClient().LRange(ctx, key, start, end).Result()
		fmt.Println(logs)
		if len(logs) > 0 {
			w := bufio.NewWriter(file)
			for _, log := range logs {
				_, writeerr := fmt.Fprintln(w, log)
				if writeerr != nil {
					continue
				}
			}
			finishErr := w.Flush()
			if finishErr == nil {
				//删除redis中已写完数据
				redisClient.GetRedisClient().LTrim(ctx, key, end+1, -1).Result()
			}
		} else {
			break
		}
	}
}

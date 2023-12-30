package db

import (
	"log"
	"strings"

	"cmsApp/configs"
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"cmsApp/pkg/redisClient"

	"github.com/spf13/cobra"
)

var cmdSeed = &cobra.Command{
	Use:   "seed [-t table]",
	Short: "DB Seed",
	Run:   seedFunc,
}

var tableSeed string
var confPath string

func init() {
	cmdSeed.Flags().StringVarP(&confPath, "config path", "c", "", "config path")
	cmdSeed.Flags().StringVarP(&tableSeed, "table", "t", "", "input a table name")
}

func seedFunc(cmd *cobra.Command, args []string) {
	var tableMap map[string]struct{}
	var err error

	err = configs.Init(configPath)
	if err != nil {
		log.Fatalf("start fail:[Config Init] %s", err.Error())
	}

	err = redisClient.Init()
	if err != nil {
		log.Fatalf("start fail:[Redis Init] %s", err.Error())
	}

	err = postgresqlx.Init()
	if err != nil {
		log.Fatalf("start fail:[Mysql Init] %s", err.Error())
	}

	tableMap = make(map[string]struct{})
	if tableSeed != "" {
		tablesSlice := strings.Split(tableSeed, ",")
		for _, v := range tablesSlice {
			tableMap[v] = struct{}{}
		}
	}

	for _, v := range models.GetModels() {

		if tableSeed != "" {
			if _, ok := tableMap[v.(postgresqlx.GaTabler).TableName()]; !ok {
				continue
			}
		}

		tabler := v.(postgresqlx.GaTabler)
		db := postgresqlx.GetDB(tabler)
		tabler.FillData(db)
	}

}

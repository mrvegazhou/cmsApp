package postgresqlx

import (
	"cmsApp/configs"
	"cmsApp/pkg/loggers"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var mapDB map[string]*gorm.DB

type GaTabler interface {
	schema.Tabler
	FillData(*gorm.DB)
	GetConnName() string
}

type BaseModle struct {
	ConnName string `gorm:"-" json:"-"`
}

func (b *BaseModle) TableName() string {
	return ""
}

func (b *BaseModle) FillData(db *gorm.DB) {}

func (b *BaseModle) GetConnName() string {
	return b.ConnName
}

// 获取链接
func GetDB(model GaTabler) *gorm.DB {

	db, ok := mapDB[model.GetConnName()]
	if !ok {
		errMsg := fmt.Sprintf("connection name%s no exists", model.GetConnName())
		loggers.LogError(context.Background(), "get_db_error", "GetDB", map[string]string{
			"msg": errMsg,
		})
	}
	return db
}

func Init() error {
	mapDB = make(map[string]*gorm.DB)

	for _, postgresConfig := range configs.App.Postgres {
		var err error
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", postgresConfig.Host, postgresConfig.Port, postgresConfig.User, postgresConfig.Password, postgresConfig.DBName)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			//Logger: logger.Discard,
			Logger: logger.New(log.New(os.Stdout, "[GORM]\u0020", log.LstdFlags), logger.Config{
				SlowThreshold:             100 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			}),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
		if err != nil {
			return err
		}
		sqlDb, _ := db.DB()
		sqlDb.SetMaxOpenConns(postgresConfig.MaxOpenConn)
		sqlDb.SetMaxIdleConns(postgresConfig.MaxIdleConn)

		//注册回调函数
		RegisterCallback(db)

		mapDB[postgresConfig.Name] = db
	}

	return nil

}

func RegisterCallback(db *gorm.DB) {
	//注册创建数据回调
	db.Callback().Create().After("gorm:create").Register("my_plugin:after_create", func(db *gorm.DB) {
		str := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		loggers.LogInfo(db.Statement.Context, "sql", "create sql", map[string]string{
			"info": str,
		})
	})
	db.Callback().Query().After("gorm:query").Register("my_plugin:after_select", func(db *gorm.DB) {
		str := fmt.Sprintf("sql语句：%s 参数：%s", db.Statement.SQL.String(), db.Statement.Vars)

		loggers.LogInfo(db.Statement.Context, "sql", "query sql", map[string]string{
			"info": str,
		})
	})
	//TODO 注册删除数据回调
	db.Callback().Delete().After("gorm:delete").Register("my_plugin:after_delete", func(db *gorm.DB) {
		str := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		loggers.LogInfo(db.Statement.Context, "sql", "delete sql", map[string]string{
			"info": str,
		})
	})
	//TODO 注册更新数据回调
	db.Callback().Update().After("gorm:update").Register("my_plugin:after_update", func(db *gorm.DB) {
		str := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		loggers.LogInfo(db.Statement.Context, "sql", "update sql", map[string]string{
			"info": str,
		})
	})
}

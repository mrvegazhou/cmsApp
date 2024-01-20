package configs

import (
	"flag"
	"os"
	"testing"

	"cmsApp/pkg/utils/filesystem"

	"gopkg.in/yaml.v2"
)

var RootPath string

type AppConf struct {
	Postgres     []PostgresConf `yaml:"postgres" json:"postgres"`
	Redis        RedisConf      `yaml:"redis" json:"redis"`
	Session      SessionConf    `yaml:"session" json:"session"`
	Base         BaseConf       `yaml:"base" json:"base"`
	Captcha      Captcha        `yaml:"captcha" json:"captcha"`
	Login        Login          `yaml:"login" json:"login"`
	Rsa          Rsa            `yaml:"rsa" json:"rsa"`
	Email        Email          `yaml:"email" json:"email"`
	Register     Register       `yaml:"register" json:"register"`
	SlideCaptcha SlideCaptcha   `yaml:"slide_captcha" json:"slide_captcha"`
	Upload       Upload         `yaml:"upload" json:"upload"`
	Article      Article        `yaml:"article" json:"article"`
}
type Article struct {
	JwtSecret string `yaml:"jwt_secret" json:"jwt_secret"`
}
type Upload struct {
	BasePath       string   `yaml:"base_path" json:"basePath"`
	ImageAllowExts []string `yaml:"image_allow_exts" json:"imageAllowExts"`
	ImageMaxSize   int      `yaml:"image_max_size" json:"imageMaxSize"`
	Key            string   `yaml:"key" json:"key"`
}

type Login struct {
	Times uint `yaml:"times" json:"times"`
	//ExpirTimeDuration int    `yaml:"expir_time_duration" json:"expir_time_duration"`
	JwtSecret        string `yaml:"jwt_secret" json:"jwt_secret"`
	JwtRefreshSecret string `yaml:"jwt_refresh_secret" json:"jwt_refresh_secret"`
	JwtExp           uint   `yaml:"jwt_exp" json:"jwt_exp"` // 过期的小时数
}

type Register struct {
	LimitRate     int `yaml:"limit_rate" json:"limit_rate"`
	LimitDuration int `yaml:"limit_duration" json:"limit_duration"`
}

type Rsa struct {
	PrivateStr string `yaml:"private_str" json:"private_str"`
	PublicStr  string `yaml:"public_str" json:"public_str"`
}

// 简单验证码配置
type Captcha struct {
	ClickTimes      int `yaml:"click_times" json:"click_times"`
	Count           int `yaml:"count" json:"count"`
	IntervalTime    int `yaml:"interval_time" json:"interval_time"`
	RecoveryTime    int `yaml:"recovery_time" json:"recovery_time"`
	IntervalKeyTime int `yaml:"interval_key_time" json:"interval_key_time"`
}

// 滑动、猜字高级验证码配置
type SlideCaptcha struct {
	CacheType      string       `yaml:"cache_type"`
	Watermark      *Watermark   `yaml:"watermark"`
	ClickWord      *ClickWord   `yaml:"click_word"`
	BlockPuzzle    *BlockPuzzle `yaml:"block_puzzle"`
	CacheExpireSec int          `yaml:"cache_expire_sec"`
	ImagesPath     *ImagesPath  `yaml:"images_path"`
	DefaultFont    string       `yaml:"default_font"`
}
type ImagesPath struct {
	DefaultTmpImgDir     string `yaml:"default_tmp_img_dir"`
	DefaultBgImgDir      string `yaml:"default_bg_img_dir"`
	DefaultClickBgImgDir string `yaml:"default_click_bg_img_dir"`
}
type Watermark struct {
	FontSize int    `yaml:"font_size"`
	Color    string `yaml:"color"`
	Text     string `yaml:"text"`
}
type ClickWord struct {
	FontSize int `yaml:"font_size"`
	FontNum  int `yaml:"font_num"`
}
type BlockPuzzle struct {
	// 校验时 容错偏移量
	Offset int `yaml:"offset"`
}

type PostgresConf struct {
	Name        string `yaml:"name" json:"name"`
	Host        string `yaml:"host" json:"host"`
	Port        string `yaml:"port" json:"port"`
	User        string `yaml:"user" json:"user"`
	Password    string `yaml:"password" json:"password"`
	DBName      string `yaml:"dbname" json:"dbname"`
	MaxOpenConn int    `yaml:"max_open_conn" json:"max_open_conn"`
	MaxIdleConn int    `yaml:"max_idle_conn" json:"max_idle_conn"`
}

type RedisConf struct {
	Addr     string `yaml:"addr" json:"addr"`
	Db       int    `yaml:"db" json:"db"`
	Password string `yaml:"password" json:"password"`
}

type SessionConf struct {
	SessionName string `yaml:"session_name"`
}

type Email struct {
	Value               string `yaml:"value" json:"value"`
	Password            string `yaml:"password" json:"password"`
	Smtp                string `yaml:"smtp" json:"smtp"`
	SmtpPort            string `yaml:"smtp_port" json:"smtp_port"`
	PoolSize            int    `yaml:"pool_size" json:"pool_size"`
	EmailName           string `yaml:"email_name" json:"email_name"`
	SendExpirDuration   int    `yaml:"send_expir_duration" json:"send_expir_duration"`
	SendEmailLimitCount int    `yaml:"send_email_limit_count" json:"send_email_limit_count"`
}

type BaseConf struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	LogMedia string `yaml:"log_media"`
	PageSize int    `yaml:"page_size"`
}

var App *AppConf

// 初始化配置文件
func Init(path string) error {

	var err error
	if path == "" {
		RootPath, err = filesystem.RootPath()
		if err != nil {
			return err
		}
	} else {
		RootPath = path
	}

	//否则执行 go test 报错
	testing.Init()
	flag.Parse()

	yamlFile, err := os.ReadFile(RootPath + "/configs/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &App)
	if err != nil {
		return err
	}
	return nil
}

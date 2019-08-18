package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	PageSize        int
	PrefixUrl       string
	RuntimeRootPath string

	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowExts  []string

	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
	ConnTimeout int
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var (
	AppSetting = &App{}
	ServerSetting = &Server{}
	DatabaseSetting = &Database{}
	RedisSetting = &Redis{}
	cfg *ini.File
)

func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fial to parse 'conf/app.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapto map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("cfg.MapTo %s setting err: %v", section, err)
	}
}
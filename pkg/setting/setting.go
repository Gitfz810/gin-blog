package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort int
	ReadTimeOut time.Duration
	WriteTimeOut time.Duration

	PageSize int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	ser, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	HTTPPort = ser.Key("HTTP_PORT").MustInt(8000)
	ReadTimeOut = time.Duration(ser.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeOut = time.Duration(ser.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	app, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = app.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = app.Key("PAGE_SIZE").MustInt(10)
}

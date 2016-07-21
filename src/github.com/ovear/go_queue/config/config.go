// config
package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

var configFile string = "config.ini"
var httpPort int
var logFile string
var authKey string
var threadNum int

var dbHost string
var dbPort string
var dbUsername string
var dbPassword string
var dbDatabase string

func init() {
	succ := InitConfig()
	if !succ {
		fmt.Println("config load failed, system exit..")
		os.Exit(-1)
	}
}

func InitConfig() (succ bool) {
	fmt.Println("loading config...")
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Println("config load failed")
		return
	}

	httpSec, err := cfg.GetSection("http")
	if err != nil {
		fmt.Println("config load failed")
		return
	}
	otherSec, err := cfg.GetSection("other")
	if err != nil {
		fmt.Println("config load failed")
		return
	}
	dbSec, err := cfg.GetSection("db")
	if err != nil {
		fmt.Println("config load failed")
		return
	}

	httpPort, err = httpSec.Key("port").Int()
	if err != nil {
		fmt.Println("config load failed")
		return
	}
	threadNum, err = otherSec.Key("threadNum").Int()
	if err != nil {
		fmt.Println("config load failed")
		return
	}
	logFile = otherSec.Key("logFile").String()
	authKey = otherSec.Key("authKey").String()

	//db config load
	dbHost = dbSec.Key("host").String()
	dbPort = dbSec.Key("port").String()
	dbUsername = dbSec.Key("username").String()
	dbPassword = dbSec.Key("password").String()
	dbDatabase = dbSec.Key("database").String()

	succ = true
	fmt.Println("load config success")
	return
}

func GetLogFile() string {
	return logFile
}

func GetHttpPort() int {
	return httpPort
}

func GetAuthKey() string {
	return authKey
}

func GetDbHost() string {
	return dbHost
}

func GetDbPort() string {
	return dbPort
}

func GetDbUsername() string {
	return dbUsername
}

func GetDbPassword() string {
	return dbPassword
}

func GetDbDatabse() string {
	return dbDatabase
}

func GetThreadNum() int {
	return threadNum
}

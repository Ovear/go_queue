//数据库处理包
package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ovear/go_queue/config"
	"github.com/ovear/go_queue/logger"
)

var host string
var port string
var username string
var password string
var database string

//连接字符串
//username:password@tcp(ip:3306)/dbname
var formatString = "%s:%s@tcp(%s:%s)/%s"

var db *sql.DB

func init() {
	InitDb()
}

//初始化数据库
func InitDb() {
	host = config.GetDbHost()
	port = config.GetDbPort()
	username = config.GetDbUsername()
	password = config.GetDbPassword()
	database = config.GetDbDatabse()

	connectStr := fmt.Sprintf(formatString, username, password, host, port, database)
	logger.InfoF("db init [%s]", connectStr)

	//数据是否有效
	mdb, err := sql.Open("mysql", connectStr)
	if err != nil {
		logger.Fatal("db init failed", err)
	}
	//测试是否能连上数据库
	err = mdb.Ping()
	if err != nil {
		logger.Fatal("db init failed", err)
	}
	mdb.SetMaxOpenConns(100)
	mdb.SetMaxIdleConns(30)
	db = mdb

	logger.Info("db init success")
}

func GetDb() *sql.DB {
	if db == nil {
		InitDb()
	}
	return db
}

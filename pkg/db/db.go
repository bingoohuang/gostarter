package db

import (
	"github.com/bingoohuang/gostarter/pkg/conf"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // 默认是用sqlite3
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

var X *sqlx.DB

func InitDB() *sqlx.DB {

	dbConf := conf.Conf.Db
	if dbConf.Driver == "" {
		logrus.Info("no driver")
		return nil
	}
	var err error
	X, err = sqlx.Open(dbConf.Driver, dbConf.DataSource)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}

	X.SetMaxIdleConns(dbConf.MaxIdleConnections)
	X.SetMaxOpenConns(dbConf.MaxOpenConnections)
	X.SetConnMaxIdleTime(time.Duration(dbConf.ConnMaxIdleTime) * time.Minute)
	X.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime) * time.Minute)

	createDemoTable()
	return X
}

func createDemoTable() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		age INTEGER DEFAULT 0
	);`

	_, err := X.Exec(schema)
	if err != nil {
		log.Fatalf("无法创建表: %v", err)
	}
}

package core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var GormDB *gorm.DB

func InitDB() {
	db, err := gorm.Open(mysql.Open(SysConfig.DBConfig.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqldb.SetMaxOpenConns(SysConfig.DBConfig.MaxOpenConn)
	sqldb.SetMaxIdleConns(SysConfig.DBConfig.MinIdleConn)
	sqldb.SetConnMaxLifetime(time.Duration(SysConfig.DBConfig.MaxLifeSecond) * time.Second)
	GormDB = db
}

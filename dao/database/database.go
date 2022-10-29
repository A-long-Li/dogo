/**
 *@filename       mysql.go
 *@Description
 *@author          liyajun
 *@create          2022-10-28 23:49
 */

package database

import (
	"fmt"
	"net/url"

	"github.com/A-long-Li/dogo/settings"
	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func Init(conf *settings.DatabaseConfig) (err error) {
	driverName := conf.DriverName
	host := conf.Host
	port := conf.Port
	database := conf.DatabaseName
	username := conf.User
	password := conf.PassWord
	charset := conf.CharSet
	loc := conf.Location

	msg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	db, err = gorm.Open(driverName, msg)
	fmt.Println(msg)
	if err != nil {
		zap.L().Error("connect database failed", zap.Error(err))
	}
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)

	return
}

func Close() {
	_ = db.Close()
}

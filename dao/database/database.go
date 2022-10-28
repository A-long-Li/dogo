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

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	db *gorm.DB
)

func Init() (err error) {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")

	msg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	db, err = gorm.Open(driverName, msg)
	if err != nil {
		zap.L().Error("connect database failed", zap.Error(err))
	}
	db.DB().SetMaxOpenConns(viper.GetInt("datasource.max_open_conn"))
	db.DB().SetMaxIdleConns(viper.GetInt("datasource.max_idle_conn"))

	return db.DB().Ping()
}

func Close() {
	_ = db.Close()
}

/**
 *@filename       main.go
 *@Description
 *@author          liyajun
 *@create          2022-10-28 23:02
 */
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dogo/dao/database"
	"dogo/dao/redis"
	"dogo/logger"
	"dogo/routes"
	"dogo/settings"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

func main() {

	initConfig()
	initLogger()
	initDatabase()
	defer database.Close()
	//initRedis()
	defer zap.L().Sync()

	//defer redis.Close()
	r := routes.SetUp()
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": settings.Conf,
		})
	})
	//平滑关机 重启项目
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}

func initConfig() {
	err := settings.Init()
	if err != nil {
		panic(fmt.Errorf("init settings failed,err:%s\n", err))
	}
}
func initLogger() {
	err := logger.Init(settings.Conf.LogConfig)
	if err != nil {
		panic(fmt.Errorf("logger init failed,err:%s\n", err))
	}
	zap.L().Debug("logger init success.....")
}
func initDatabase() {
	err := database.Init(settings.Conf.DatabaseConfig)
	if err != nil {
		panic(fmt.Errorf("database -> %s init failed,err:%s\n", settings.Conf.DatabaseConfig.DriverName, err))
	}
}

func initRedis() {
	err := redis.Init(settings.Conf.RedisConfig)

	if err != nil {
		panic(fmt.Errorf("database -> %s init failed,err:%s\n", viper.GetString("datasource.driverName"), err))
	}
}

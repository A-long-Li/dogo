/**
 *@filename       redis.go
 *@Description
 *@author          liyajun
 *@create          2022-10-29 0:04
 */

package redis

import (
	"context"
	"fmt"
	"time"
	"web_app/settings"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

// 初始化连接

func Init(conf *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.PassWord, // no password set
		DB:       conf.DB,       // use default DB
		PoolSize: conf.PoolSize, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}
func Close() {
	_ = rdb.Close()
}

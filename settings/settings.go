/**
 *@filename       settings.go
 *@Description
 *@author          liyajun
 *@create          2022-10-28 23:07
 */

package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigFile("./settings/config.yaml") // 指定配置文件路径
	err = viper.ReadInConfig()                    // 查找并读取配置文件
	if err != nil {                               // 处理读取配置文件的错误
		panic(fmt.Errorf("viper.ReadInConfig() failed: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file changed %s \n:", e.Name)
	})
	return
}

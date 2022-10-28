/**
 *@filename       routes.go
 *@Description
 *@author          liyajun
 *@create          2022-10-29 0:15
 */

package routes

import (
	"web_app/logger"

	"github.com/gin-gonic/gin"
)

func SetUp() (r *gin.Engine) {
	r = gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	return
}

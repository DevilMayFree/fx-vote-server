package initialize

import (
	"fx-vote-server/common/app"
	"fx-vote-server/common/constant"
	"fx-vote-server/common/ginzap"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	gin.SetMode(getGinMode())
	engine := gin.New()

	var handlers []gin.HandlerFunc
	ginPlugins(&handlers)

	engine.Use(handlers...)

	return engine
}

// ginPlugins 设置gin插件
func ginPlugins(p *[]gin.HandlerFunc) {
	if app.Config.Gin.LogZap {
		*p = append(*p, ginzap.Logger())
	}

	ginzap.LogRequest = &app.Config.Gin.LogRequest
	ginzap.LogStack = &app.Config.Gin.LogStack
	*p = append(*p, ginzap.Recovery())
}

// getGinMode 获取运行GinMode
func getGinMode() string {
	switch app.Env {
	case constant.DEV:
		return gin.DebugMode
	case constant.TEST:
		return gin.TestMode
	case constant.PROD:
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}

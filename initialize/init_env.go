package initialize

import (
	"fx-vote-server/common/app"
	"go.uber.org/zap"
)

// InitEnv 初始化项目
func InitEnv() {
	Config()
	app.Log = Zap()
	zap.ReplaceGlobals(app.Log)
	PrintConfig()
	app.Redis = Redis()
	runServer()
}

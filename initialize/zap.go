package initialize

import (
	"fmt"
	"fx-vote-server/common/app"
	"fx-vote-server/common/gzap"
	"fx-vote-server/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap() (log *zap.Logger) {

	director := app.Config.Zap.Director
	if ok, _ := utils.PathExists(director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", director)
		_ = os.Mkdir(director, os.ModePerm)
	}

	cores := gzap.Zap.GetZapCores()
	log = zap.New(zapcore.NewTee(cores...))

	app.SqlLog = log
	app.GinLog = log
	if app.Config.Zap.ShowLine {
		log = log.WithOptions(zap.AddCaller())
	}
	return log
}

package gzap

import (
	"fx-vote-server/common/app"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

var Rotatelogs = new(_rotatelogs)

type _rotatelogs struct {
}

// GetLogFileWriteSyncer 输出到日志文件WriteSyncer
func (r *_rotatelogs) GetLogFileWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(app.Config.Zap.Director, "%Y-%m-%d", level+".log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(app.Config.Zap.MaxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return zapcore.AddSync(fileWriter), err
}

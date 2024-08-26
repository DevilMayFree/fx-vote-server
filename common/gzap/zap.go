package gzap

import (
	"fmt"
	"fx-vote-server/common/app"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap = new(_zap)

type _zap struct {
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core
func (z *_zap) GetZapCores() []zapcore.Core {

	cores := make([]zapcore.Core, 0, 7)
	l := app.Config.Zap.GetLevel()
	for ; l <= zapcore.FatalLevel; l++ {
		cores = append(cores, z.GetEncoderCore(l))
	}

	return cores
}

// GetEncoderCore 获取zapcore encoderCore
func (z *_zap) GetEncoderCore(l zapcore.Level) zapcore.Core {
	priority := z.GetLevelPriority(l)
	writer := z.GetWriteSyncer(l.String())

	cfg := z.GetEncoderConfig()
	encoder := app.Config.Zap.GetEncoder(&cfg)

	return zapcore.NewCore(encoder, writer, priority)
}

// GetWriteSyncer 获取zapcore.writeSyncer
func (z *_zap) GetWriteSyncer(level string) zapcore.WriteSyncer {

	syncers := make([]zapcore.WriteSyncer, 0, 2)

	if app.Config.Zap.LogInConsole {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if app.Config.Zap.LogInFile {
		writer, err := Rotatelogs.GetLogFileWriteSyncer(level)
		if err != nil {
			fmt.Printf("Get LogFile Write Syncer Failed err:%v", err.Error())
		} else {
			syncers = append(syncers, zapcore.AddSync(writer))
		}
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

// GetEncoderConfig 获取zapcore.EncoderConfig
func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:       "message",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    app.Config.Zap.StacktraceKey,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      app.Config.Zap.GetLevelEncoder(), // zapcore.CapitalColorLevelEncoder, // env.Config.Zap.ZapEncodeLevel(),
		EncodeTime:       z.CustomTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     z.CustomCallerEncoder,
		ConsoleSeparator: "  ",
	}
}

// CustomTimeEncoder 自定义日志前缀与输出时间格式
func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(app.Config.Zap.Prefix)
	encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

// CustomCallerEncoder 自定义日志显示全路径还是包文件
func (z *_zap) CustomCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if app.Config.Zap.FullCaller {
		zapcore.FullCallerEncoder(caller, enc) // full/path/to/package/file
	} else {
		zapcore.ShortCallerEncoder(caller, enc) // package/file
	}
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
func (z *_zap) GetLevelPriority(zapLevel zapcore.Level) zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == zapLevel
	}
}

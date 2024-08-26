package configs

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	JSONFormat = "json"

	DebugLevel  = "debug"
	InfoLevel   = "info"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DPanicLevel = "dpanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
)

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 日志级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 日志格式化  json | console
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志输出文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级 capital | capitalColor | lowerCaseColor | lowercase
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名

	MaxAge       int  `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间
	ShowLine     bool `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	FullCaller   bool `mapstructure:"full-caller" json:"full-caller" yaml:"full-caller"`          // 显示完整调用行
	LogInConsole bool `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 日志是否输出到控制台
	LogInFile    bool `mapstructure:"log-in-file" json:"log-in-file" yaml:"log-in-file"`          // 日志是否输出到文件
}

// GetEncoder 根据配日志和EncoderConfig获取Encoder
func (z *Zap) GetEncoder(cfg *zapcore.EncoderConfig) zapcore.Encoder {
	f := strings.ToLower(z.Format)
	if JSONFormat == f {
		return zapcore.NewJSONEncoder(*cfg) // Format: json
	}
	return zapcore.NewConsoleEncoder(*cfg) // Format: console
}

// GetLevel 根据配置转化为 zapcore.Level
func (z *Zap) GetLevel() zapcore.Level {
	switch strings.ToLower(z.Level) {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case DPanicLevel:
		return zapcore.DPanicLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// GetLevelEncoder 根据配置转为 zapcore.LevelEncoder
func (z *Zap) GetLevelEncoder() zapcore.LevelEncoder {
	switch strings.ToLower(z.EncodeLevel) {
	case "capital":
		return zapcore.CapitalLevelEncoder
	case "capitalcolor":
		return zapcore.CapitalColorLevelEncoder
	case "lowercasecolor":
		return zapcore.LowercaseColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// IsColorPrint 是否开启彩色打印
func (z *Zap) IsColorPrint() bool {
	l := strings.ToLower(z.EncodeLevel)
	return strings.Index(l, "color") != -1
}

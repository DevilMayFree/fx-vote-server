package configs

type Gin struct {
	LogZap     bool `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`             // ginLog是否使用zap写日志
	LogRequest bool `mapstructure:"log-request" json:"log-request" yaml:"log-request"` // recovery是否记录panic时request信息
	LogStack   bool `mapstructure:"log-stack" json:"log-stack" yaml:"log-stack"`       // recovery是否记录panic时堆栈信息
}

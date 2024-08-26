package configs

type YamlConfig struct {
	APP    APP    `mapstructure:"app" json:"app" yaml:"app"`
	Server Server `mapstructure:"server" json:"server" yaml:"server"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Gin    Gin    `mapstructure:"gin" json:"gin" yaml:"gin"`
}

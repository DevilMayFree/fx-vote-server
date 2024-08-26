package configs

type Server struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"` // 服务端口
	Name string `mapstructure:"name" json:"name" yaml:"name"` // 服务名
}

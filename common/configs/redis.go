package configs

import "runtime"

type Redis struct {
	Address      string `mapstructure:"address" json:"address" yaml:"address"`                      // 连接地址:端口
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // Redis >= 6.0 ACL控制
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 密码
	DB           int    `mapstructure:"db" json:"db" yaml:"db"`                                     // 指定数据库
	PoolSize     int    `mapstructure:"pool-size" json:"pool-size" yaml:"pool-size"`                // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
	MinIdleConns int    `mapstructure:"min-idle-conns" json:"min-idle-conns" yaml:"min-idle-conns"` // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	DialTimeout  int    `mapstructure:"dial-timeout" json:"dial-timeout" yaml:"dial-timeout"`       // 连接建立超时时间，默认5秒
	ReadTimeout  int    `mapstructure:"read-timeout" json:"read-timeout" yaml:"read-timeout"`       // 读超时，默认3秒， -1表示取消读超时
	WriteTimeout int    `mapstructure:"write-timeout" json:"write-timeout" yaml:"write-timeout"`    // 写超时，默认等于读超时
}

func (r *Redis) GetPoolSize() int {
	if r.PoolSize == 0 {
		return 4 * runtime.GOMAXPROCS(runtime.NumCPU())
	}
	return r.PoolSize
}

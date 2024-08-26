package constant

import "time"

const (
	DEV  = "dev"
	TEST = "test"
	PROD = "prod"

	ConfigType       = "yaml"
	ConfigPathFormat = "configs/config-%s.yaml"

	RedisVoteKey = "vote:key"

	RedisIpKey = "vote:ip_set"

	RedisIpExpiration = 24 * time.Hour
)

var TeacherList = []string{"陳志強", "張文輝", "杜文瀚", "黃志偉", "林建宏", "蔣鴻志", "李國峰"}

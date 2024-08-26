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

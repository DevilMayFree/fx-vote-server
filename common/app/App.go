package app

import (
	"fx-vote-server/common/configs"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Env string

	Config *configs.YamlConfig

	Log    *zap.Logger
	SqlLog *zap.Logger
	GinLog *zap.Logger

	Redis *redis.Client
)

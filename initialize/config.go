package initialize

import (
	"flag"
	"fmt"
	"fx-vote-server/common/app"
	"fx-vote-server/common/constant"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var profile string
var ConfigFile string

func Config() {
	flag.StringVar(&profile, "c", "", "choose config mode. dev | test | prod")
	flag.Parse()

	runProfile := constant.DEV // 默认开发环境
	if profile != "" {
		runProfile = strings.ToLower(profile)
	}
	app.Env = runProfile
	ConfigFile = fmt.Sprintf(constant.ConfigPathFormat, runProfile)

	v := viper.New()
	v.SetConfigFile(ConfigFile)
	v.SetConfigType(constant.ConfigType)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config File Changed:", in.Name)
		if err := v.Unmarshal(&app.Config); err != nil {
			_ = fmt.Errorf("Fatal error change config file: %s \n", err)
		}
	})

	if err := v.Unmarshal(&app.Config); err != nil {
		_ = fmt.Errorf("Fatal error change config file: %s \n", err)
	}
}

func PrintConfig() {
	app.GinLog.Info("Using Config File From :" + ConfigFile)
}

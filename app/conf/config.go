package conf

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Settings 全局设置
var Settings *config

type config struct {
	Storage   AliYun `json:"Storage" yaml:"Storage"`
	LocalPath string `json:"LocalPath" yaml:"LocalPath"`
}

type AliYun struct {
	AccessKeyID     string
	AccessKeySecret string
	EndPoint        string
	DefaultBucket   string
	Path            string
}

// SetupSetting Setup initializes the configuration instance
func SetupSetting() error {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "missing config.yaml")
	} // Find and read the config file
	viper.WatchConfig()
	var t *time.Timer
	// 热重载配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 此处监听文件修改， 为防止响应过快多次加载，等1m之后再重新加载文件
		if t != nil {
			t.Stop()
		}
		t = time.AfterFunc(1*time.Minute, func() {
			if err := loadConfig(); err != nil {
				log.Printf("viper onConfigChange failed %v", err)
			}
		})
	})
	return loadConfig()
}

func loadConfig() error {
	Settings = &config{}
	if err := viper.Unmarshal(Settings); err != nil {
		return errors.Wrap(err, "parse config.yaml failed")
	}
	return nil
}

package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// CF -> for use configs model
	CF = &Configs{}
)

// Configs config models
type Configs struct {
	UniversalTranslator *ut.UniversalTranslator
	Validator           *validator.Validate
	App                 struct {
		Port int    `mapstructure:"PORT"`
		URL  string `mapstructure:"URL"`
	} `mapstructure:"APP"`
	SQL struct {
		MySQL struct {
			Host         string `mapstructure:"HOST"`
			Port         int    `mapstructure:"PORT"`
			Username     string `mapstructure:"USERNAME"`
			Password     string `mapstructure:"PASSWORD"`
			DatabaseName string `mapstructure:"DATABASE_NAME"`
		} `mapstructure:"MY_SQL"`
		PostgreSQL struct {
			Host         string `mapstructure:"HOST"`
			Port         int    `mapstructure:"PORT"`
			Username     string `mapstructure:"USERNAME"`
			Password     string `mapstructure:"PASSWORD"`
			DatabaseName string `mapstructure:"DATABASE_NAME"`
		} `mapstructure:"POSTGRE_SQL"`
	} `mapstructure:"SQL"`
	Mongo struct {
		Host         string `mapstructure:"HOST"`
		Port         int    `mapstructure:"PORT"`
		Username     string `mapstructure:"USERNAME"`
		Password     string `mapstructure:"PASSWORD"`
		DatabaseName string `mapstructure:"DATABASE_NAME"`
	} `mapstructure:"MONGO"`
}

// InitConfig init config
func InitConfig(configPath string, environment string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(fmt.Sprintf("config.%s", environment))
	v.AutomaticEnv()
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := bindingConfig(v, CF); err != nil {
			logrus.Error("binding error:", err)
			return
		}
	})

	return nil
}

// bindingConfig binding config
func bindingConfig(vp *viper.Viper, cf *Configs) error {
	if err := vp.Unmarshal(&cf); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}

	validate := validator.New()
	cf.Validator = validate

	return nil
}

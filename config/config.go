package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type Config struct {
	DBHost         string `mapstructure:"db_host"`
	DBPort         int    `envconfig:"db_port"`
	DBUsername     string `mapstructure:"db_username"`
	DBPassword     string `mapstructure:"db_password"`
	DBDatabase     string `envconfig:"db_database"`
	DBReadTimeout  string `envconfig:"db_read_timeout"`
	DBMaxOpenConns int    `envconfig:"db_max_open_conns"`
	DBMaxIdleConns int    `envconfig:"db_max_idle_conns"`
}

func Init(filename string) Config {
	var Conf Config
	t := reflect.TypeOf(Conf)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if env, ok := f.Tag.Lookup("mapstructure"); ok {
			_ = viper.BindEnv(env)
		}
	}

	viper.SetConfigFile(filename)
	viper.SetConfigType("env")

	_ = viper.ReadInConfig()
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Sprintf("could not unmarshal config: %s", err.Error()))
	}
	return Conf
}

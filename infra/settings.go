package infra

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetEnv(v interface{}) error {
	viper.SetConfigFile((".env"))
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found error")
		} else {
			return fmt.Errorf("config file has error : %s ", err.Error())
		}
	}

	if err := viper.Unmarshal(&v); err != nil {
		return fmt.Errorf("parse env error : %s", err)
	}

	fmt.Println(v)

	return nil
}

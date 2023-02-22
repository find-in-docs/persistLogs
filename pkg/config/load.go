package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load() {

  /*
  pgLoginUserFilename := "/mnt/secret/pg-login/user-value-path"
  pgLoginPassFilename := "/mnt/secret/pg-login/pass-value-path"
  */

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/mnt/config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file was not found.\n")
		} else {
			panic(fmt.Errorf("Fatal error reading config file: error: %w\n", err))
		}
	}
}

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
func loadFile(fn string) ([]byte, error) {
  b, err := os.ReadFile(fn)
  if err != nil {
    fmt.Printf("Could not read file: %s, err: %s\n", fn, err)
    return nil, err
  }

  return b, nil
}

func loadString(fn string) (string, error) {
  b, err := loadFile(fn)
  if err != nil {
    return nil, err
  }

  return b, nil
}
*/

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

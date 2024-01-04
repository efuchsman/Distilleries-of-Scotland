package distilleries

import (
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")
	configPath := filepath.Join(projectRoot, "config", "config_test.yml")

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}
}

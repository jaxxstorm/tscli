package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	fileName = ".tscli"
	fileType = "yaml"
)

func Init() {
	v := viper.GetViper()

	// Search order: cwd ⇒ $HOME
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(".")
	if home, _ := os.UserHomeDir(); home != "" {
		v.AddConfigPath(home)
	}

	_ = v.ReadInConfig() // ignore “not found”
	v.SetDefault("output", "json")
}

func Save() error {
	v := viper.GetViper()
	path := v.ConfigFileUsed()
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, fileName+"."+fileType)
	}
	return v.WriteConfigAs(path)
}

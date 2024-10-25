package kqueuey

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct{}

type BadgerDB struct {
	dirs            []string
	numCompactors   int
	compressionType string
	syncWrites      bool
}

type Raft struct {
	addrs   []string
	nodeIds []string
}

func GetViper() {
	viper.SetConfigName("kqueuey-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println(viper.Get("badgerDB"))
	fmt.Println(viper.Get("raft"))
}

func setDefaults() {}

func setEnv() {}

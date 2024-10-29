package kqueuey

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

var (
	// Todo revise this.
	errConfigFileNotLocated   = errors.New("config not found, set path to look for config file: default /usr/local/etc/")
	errDuplicateStoragePath   = errors.New("storage path must be unique per node to avoid file locking conflicts")
	errDuplicateNodeId        = errors.New("node id must be unique per node to conform to rafts consensus protocol")
	errDuplicateNodePort      = errors.New("node port most be unique to allow proper communication between nodes")
	errNodeIdNotFound         = errors.New("node id not found")
	errNodeStorageDirNotFound = errors.New("node storage directory not found")
	errUnknownNodeAddress     = errors.New("unknown node address")
	errRaftClusterIdNotFound  = errors.New("cluster id was not found")
)

const (
	ConfigType     = "yaml"
	ConfigFileName = "kqueuey-config"
)

type Configuration struct {
	BadgerOpts Storage `mapstructure:"storage"`
	RaftOpts   Raft    `mapstructure:"raft"`
}

type Storage struct {
	NumCompactors   int    `mapstructure:"num_compactors"`
	CompressionType string `mapstructure:"compression_type"`
	SyncWrites      bool   `mapstructure:"sync_writes"`
}

type Raft struct {
	ClusterId string     `mapstructure:"cluster_id"`
	Nodes     []RaftNode `mapstructure:"nodes"`
}

type RaftNode struct {
	Id         string `mapstructure:"id"`
	BindAddr   string `mapstructure:"bind_addr"`
	StorageDir string `mapstructure:"storage_dir"`
}

func LoadConfiguration(flagOpts FlagOpts, logger *slog.Logger) (*Configuration, error) {
	v := initializeViper(flagOpts)

	if err := v.ReadInConfig(); err != nil {
		return nil, errConfigFileNotLocated
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.validate(logger); err != nil {
		return nil, err
	}

	return &config, nil
}

func initializeViper(flagOpts FlagOpts) *viper.Viper {
	v := viper.New()
	setDefaults(v)
	setBinds(v)

	v.AutomaticEnv()
	v.SetConfigType(ConfigType)
	v.SetConfigName(ConfigFileName)

	// Define multiple search options for getting the configuration, from directory path.
	envVariableConfigPath := v.GetString("CONFIG_PATH")
	v.AddConfigPath(envVariableConfigPath)
	v.AddConfigPath(flagOpts.ConfigPath)

	return v
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("storage.num_compactors", 4)
	v.SetDefault("storage.compression_type", CompressionSnappy)
	// Allow setting multiple options to get a config path to read config yaml file. Can be from flags or ENV variables.
	// This is for setting default location to look for config option if not found in ENV.
	v.SetDefault("CONFIG_PATH", "/usr/local/etc/")
}

func setBinds(v *viper.Viper) {
	v.BindEnv("CONFIG_PATH")
}

// Validate ensures that the config file is set up correctly, to allow proper start up of queue server.
func (c *Configuration) validate(logger *slog.Logger) error {
	if c.BadgerOpts.NumCompactors < 4 {
		logger.Warn("badger's compactors was below minimum required amount of 4. Defaulting to 4")
		c.BadgerOpts.NumCompactors = 4
	}

	if err := c.RaftOpts.validateRaftConfigOptions(); err != nil {
		return err
	}

	return nil
}

func (r *Raft) validateRaftConfigOptions() error {
	if r.ClusterId == "" {
		return errRaftClusterIdNotFound
	}

	duplicates := make(map[string]uint8)
	for _, node := range r.Nodes {
		if node.Id == "" {
			return errNodeIdNotFound
		}

		if node.StorageDir == "" {
			return errNodeStorageDirNotFound
		}

		duplicates[node.StorageDir] += 1
		if duplicates[node.StorageDir] >= 2 {
			return errDuplicateStoragePath
		}

		duplicates[node.Id] += 1
		if duplicates[node.Id] >= 2 {
			return errDuplicateNodeId
		}

		splitAddr := strings.Split(node.BindAddr, ":")
		port := splitAddr[1]
		if len(splitAddr) != 2 || len(port) < 4 {
			return errUnknownNodeAddress
		}

		duplicates[port] += 1
		if duplicates[port] == 2 {
			return errDuplicateNodePort
		}
	}

	return nil
}

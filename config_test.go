package kqueuey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig(t *testing.T) {
	t.Parallel()

	t.Run("read and parse kqueuey config with all correct values ", func(t *testing.T) {
		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "2",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path1/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.NoError(t, err)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate node id", func(t *testing.T) {
		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path1/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}

		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errDuplicateNodeId)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate node port", func(t *testing.T) {

		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "2",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path1/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errDuplicateNodePort)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate node port", func(t *testing.T) {

		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errDuplicateStoragePath)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate cluster id", func(t *testing.T) {

		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errRaftClusterIdNotFound)
	})

	t.Run("kqueuey configuration file validation failed due to unkown node address", func(t *testing.T) {
		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "2",
					BindAddr:   "127.0.0.1:",
					StorageDir: "path1/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errUnknownNodeAddress)
	})

	t.Run("kqueuey configuration file validation failed due to missing node id", func(t *testing.T) {

		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path0/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errNodeIdNotFound)
	})

	t.Run("kqueuey configuration file validation failed due missing storage directory for node", func(t *testing.T) {
		flagOpts := FlagOpts{
			ConfigPath: "",
			Logging: Logging{
				Level:  "debug",
				Format: "json",
			},
		}
		logger := flagOpts.Logging.NewLogger()
		expectedStorageConfig := Storage{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := Raft{
			ClusterId: "1000",
			Nodes: []RaftNode{
				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "",
				},

				{
					Id:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path0/dir",
				},

				{
					Id:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		config := Configuration{
			BadgerOpts: expectedStorageConfig,
			RaftOpts:   expectedRaftConfig,
		}
		err := config.validate(logger)
		assert.ErrorIs(t, err, errNodeStorageDirNotFound)
	})
}

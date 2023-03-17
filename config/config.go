package config

import "github.com/pitabwire/frame"

type PartitionConfig struct {
	*frame.ConfigurationDefault

	NotificationServiceURI string `default:"127.0.0.1:7020" envconfig:"NOTIFICATION_SERVICE_URI"`
	QueuePartitionSyncURL  string `default:"mem://partition_sync_hydra" envconfig:"QUEUE_PARTITION_SYNC"`
	PartitionSyncName      string `default:"partition_sync_hydra" envconfig:"QUEUE_PARTITION_SYNC_NAME"`
}

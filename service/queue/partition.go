package queue

import (
	"context"
	"encoding/json"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
)

type PartitionSyncQueueHandler struct {
	Service    *frame.Service

}

func (psq *PartitionSyncQueueHandler) Handle(ctx context.Context, payload []byte) error {

	partition := &models.Partition{}
	err := json.Unmarshal(payload, partition)
	if err != nil {
		return err
	}

	return business.SyncPartitionOnHydra(ctx, partition)

}
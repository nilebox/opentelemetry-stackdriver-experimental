package metrics

import (
	"context"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type Exporter struct {
	client *monitoring.MetricClient
}

var _ export.Exporter = &Exporter{}

func (e *Exporter) Export(ctx context.Context, checkpoints export.CheckpointSet) error {
	err := e.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{})
	if err != nil {
		return err
	}
	return nil
}

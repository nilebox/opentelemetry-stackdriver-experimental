package metrics

import (
	"context"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"

	metric "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
)

type Exporter struct {
	client      *monitoring.MetricClient
	projectName string
}

var _ export.Exporter = &Exporter{}

func (e *Exporter) Export(ctx context.Context, checkpoints export.CheckpointSet) error {
	var timeSeries []*monitoringpb.TimeSeries
	checkpoints.ForEach(func(record export.Record) {
		timeSeries = append(timeSeries, &monitoringpb.TimeSeries{
			Resource: &monitoredres.MonitoredResource{},
			Metric:   &metric.Metric{},
		})
	})
	err := e.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name:       e.projectName,
		TimeSeries: timeSeries,
	})
	if err != nil {
		return err
	}
	return nil
}

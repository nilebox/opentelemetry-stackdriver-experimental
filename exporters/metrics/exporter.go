package metrics

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"path"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

type Exporter struct {
	client      *monitoring.MetricClient
	projectName string
}

var _ export.Exporter = &Exporter{}

func (e *Exporter) Export(ctx context.Context, checkpoints export.CheckpointSet) error {
	var timeSeries []*monitoringpb.TimeSeries
	checkpoints.ForEach(func(record export.Record) {
		labels, err := metricLabels(record.Labels())
		if err != nil {
			panic(err)
		}
		pt, err := point(record.Aggregator())
		if err != nil {
			panic(err)
		}
		timeSeries = append(timeSeries, &monitoringpb.TimeSeries{
			Metric:   &metric.Metric{
				Type: e.metricType(record.Descriptor().Name()),
				Labels: labels,
			},
			Resource: &monitoredres.MonitoredResource{
				// TODO
			},
			Points: []*monitoringpb.Point{
				pt,
			},
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

func (e *Exporter) metricType(name string) string {
	return path.Join("custom.googleapis.com", "opentelemetry", name)
}

func metricLabels(exportLabels export.Labels) (map[string]string, error) {
	labels := make(map[string]string)
	for _, kv := range exportLabels.Ordered() {
		if _, ok := labels[string(kv.Key)]; ok {
			return nil, errors.New("duplicate label keys not allowed")
		}
		labels[string(kv.Key)] = kv.Value.AsString()
	}
	return labels, nil
}

func point(agg export.Aggregator) (*monitoringpb.Point, error) {
	if sum, ok := agg.(aggregator.Sum); ok {
		num, err := sum.Sum()
		if err != nil {
			return nil, err
		}
		return &monitoringpb.Point{
			Value: &monitoringpb.TypedValue{
				Value: &monitoringpb.TypedValue_Int64Value{
					Int64Value: num.AsInt64(),
				},
			},
		}, nil
	}
	return nil, errors.New("unknown aggregatot")
}

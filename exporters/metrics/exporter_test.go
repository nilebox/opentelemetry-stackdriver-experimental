package metrics

import (
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

var _ export.Exporter = &Exporter{}

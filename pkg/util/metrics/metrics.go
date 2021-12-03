package metrics

import (
	"fmt"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"go.opencensus.io/stats/view"
)

func Start(projectID, env, app string) (func(), error) {
	const interval = time.Minute
	getMetricPrefix := func(s string) string {
		return fmt.Sprintf("com.im-barry/metrics/%s", s)
	}
	opts := stackdriver.Options{
		ProjectID:               projectID,
		OnError:                 func(err error) {},
		GetMetricPrefix:         getMetricPrefix,
		DefaultMonitoringLabels: newLabels(env, app),
		MonitoredResource:       monitoredresource.Autodetect(),
		ReportingInterval:       interval,
	}
	exporter, err := stackdriver.NewExporter(opts)
	if err != nil {
		return nil, err
	}
	cleanup := exporter.Flush

	if err := exporter.StartMetricsExporter(); err != nil {
		return cleanup, err
	}
	cleanup = func() {
		exporter.StopMetricsExporter()
		exporter.Flush()
	}
	return cleanup, nil
}

func newLabels(env, app string) *stackdriver.Labels {
	labels := &stackdriver.Labels{}
	labels.Set("env", env, "The environment to run app")
	labels.Set("app", app, "The name of application")
	return labels
}

func NewLatencyDistribution() *view.Aggregation {
	const begin, count, exp = 25.0, 9, 2
	dist := make([]float64, count)
	dist[0] = begin
	for i := 1; i < count; i++ {
		dist[i] = dist[i-1] * exp
	}
	return view.Distribution(dist...)
}

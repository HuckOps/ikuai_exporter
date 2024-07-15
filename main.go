package main

import "ikuai_exporter/metrics"

func main() {
	metricClient := metrics.NewPrometheus()
	metricClient.NewMetrics("ikuai")
	metricClient.Run()
}

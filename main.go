package main

import (
	"flag"
	"fmt"
	"github.com/HuckOps/ikuai_exporter/metrics"
)

func main() {
	var (
		routeAddress string
		userName     string
		password     string
	)
	flag.StringVar(&routeAddress, "i", "192.168.1.1", "Route address")
	flag.StringVar(&userName, "u", "admin", "Route username")
	flag.StringVar(&password, "p", "admin", "Route password")
	flag.Parse()
	if flag.NArg() == 0 && (flag.Lookup("h") != nil || flag.Lookup("help") != nil) {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	metricClient := metrics.NewPrometheus()
	metricClient.NewMetrics("ikuai")
	metricClient.Run(routeAddress, userName, password)
}

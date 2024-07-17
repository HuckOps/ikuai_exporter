package metrics

import (
	"github.com/HuckOps/ikuai_exporter/ikuai"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

type Prometheus struct {
	Registry *prometheus.Registry
	Metrics  map[string]interface{}
}

func NewPrometheus() *Prometheus {
	registry := prometheus.NewRegistry()
	http.Handle("/metrics", promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{Registry: registry}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Println("运行exporter错误", err)
			os.Exit(1)
		}
	})
	return &Prometheus{
		Registry: registry,
	}
}

func (c *Prometheus) NewMetrics(namespace string) {
	m := MakeMetricsMap("ikuai")
	for _, metric := range m {
		if err := c.Registry.Register(metric.(prometheus.Collector)); err != nil {
			log.Fatal(err)
		}
	}
	c.Metrics = m
}

func (c *Prometheus) Run(ip, username, password string) {
	ikuaiClient := ikuai.NewClient(ip, username, password)
	ikuaiClient.Login()
	go func() {
		http.ListenAndServe(":9100", nil)
	}()
	for range time.Tick(time.Second) {
		go func() {
			sysstat := ikuaiClient.GetSysstat()
			c.Metrics["cpu"].(prometheus.Gauge).Set(sysstat.CPUPercent)
			c.Metrics["mem_buffer"].(prometheus.Gauge).Set(sysstat.Buffer)
			c.Metrics["mem_cached"].(prometheus.Gauge).Set(sysstat.Cache)
			c.Metrics["mem_total"].(prometheus.Gauge).Set(sysstat.Total)
			c.Metrics["mem_free"].(prometheus.Gauge).Set(sysstat.Free)
			c.Metrics["mem_used"].(prometheus.Gauge).Set(sysstat.MemoryUsage)
		}()
		go func() {
			ifaceStatus := ikuaiClient.GetIface()
			for _, iface := range ifaceStatus.IfaceStream {
				resultTmp := map[string]float64{
					"tx_bytes_speed": float64(iface.Upload),
					"rx_bytes_speed": float64(iface.Download),
					"tx_bytes_total": float64(iface.TotalUp),
					"rx_bytes_total": float64(iface.TotalDown),
				}
				for key, value := range resultTmp {
					c.Metrics[key].(*prometheus.GaugeVec).With(
						prometheus.Labels{"interface": iface.Interface,
							"ip_add": iface.IpAddr}).Set(value)
				}
			}
			for _, iface := range ifaceStatus.IfaceCheck {
				status := 0
				if iface.Result == "success" {
					status = 1
				}
				c.Metrics["up_link_status"].(*prometheus.GaugeVec).With(
					prometheus.Labels{"interface": iface.Interface,
						"ip_add": iface.IpAddr}).Set(float64(status))
			}
		}()
		go func() {
			lanIPItems := ikuaiClient.GetLanIPs()
			for _, ip := range lanIPItems {
				resultTmp := map[string]float64{
					"lan_device_connect_num":    float64(ip.ConnectNum),
					"lan_device_tx_bytes_speed": float64(ip.Upload),
					"lan_device_rx_bytes_speed": float64(ip.Download),
					"lan_device_tx_bytes_total": float64(ip.TotalUp),
					"lan_device_rx_bytes_total": float64(ip.TotalDown),
				}
				for key, value := range resultTmp {
					c.Metrics[key].(*prometheus.GaugeVec).With(
						prometheus.Labels{"ip_add": ip.IPAddr, "mac_addr": ip.MAC}).Set(value)
				}
			}
		}()
	}
}

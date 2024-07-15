package metrics

import "github.com/prometheus/client_golang/prometheus"

func MakeMetricsMap(namespace string) map[string]interface{} {
	var metrics map[string]interface{} = map[string]interface{}{
		"cpu": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "cpu_usage",
			Help:      "CPU usage",
			Namespace: namespace,
		}),
		"mem_buffer": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "mem_buffer",
			Help:      "Memory buffer",
			Namespace: namespace,
		}),
		"mem_cached": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "mem_cached",
			Help:      "Memory cached",
			Namespace: namespace,
		}),
		"mem_total": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "mem_total",
			Help:      "Memory total",
			Namespace: namespace,
		}),
		"mem_free": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "mem_free",
			Namespace: namespace,
		}),
		"mem_used": prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "mem_used",
			Help:      "Memory used",
			Namespace: namespace,
		}),
		"tx_bytes_speed": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "tx_bytes_speed",
			Help:      "Transmitted bytes speed",
			Namespace: namespace,
		}, []string{"interface", "ip_add"}),
		"rx_bytes_speed": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "rx_bytes_speed",
			Help:      "Received bytes speed",
			Namespace: namespace,
		}, []string{"interface", "ip_add"}),
		"tx_bytes_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "tx_bytes_total",
			Help:      "Transmitted bytes total",
			Namespace: namespace,
		}, []string{"interface", "ip_add"}),
		"rx_bytes_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "rx_bytes_total",
			Help:      "Received bytes total",
			Namespace: namespace,
		}, []string{"interface", "ip_add"}),
		"up_link_status": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "up_link_status",
			Help:      "Up link status",
			Namespace: namespace,
		}, []string{"interface", "ip_add"}),
		// lan单机
		"lan_device_connect_num": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_connect_num",
			Help:      "Lan device connect num",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
		"lan_device_tx_bytes_speed": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_tx_bytes_speed",
			Help:      "Lan device transmit bytes speed",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
		"lan_device_rx_bytes_speed": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_rx_bytes_speed",
			Help:      "Lan device receive bytes speed",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
		"lan_device_tx_bytes_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_tx_bytes_total",
			Help:      "Lan device transmit bytes total",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
		"lan_device_rx_bytes_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_rx_bytes_total",
			Help:      "Lan device receive bytes total",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
		"lan_device_up_link_status": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "lan_device_up_link_status",
			Help:      "Lan device up link status",
			Namespace: namespace,
		}, []string{"mac_addr", "ip_add"}),
	}
	return metrics
}

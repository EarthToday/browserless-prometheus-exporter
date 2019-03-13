package main

import (
	"encoding/json"
	"net/http"
	"log"
	"bytes"
	"io"
	"fmt"
	"flag"
)

type MetricsObject struct {
	Cpu float64 `json:"cpu"`
	Date int64 `json:"date"`
	Error int `json:"error"`
	Memory float64 `json:"memory"`
	Queued int `json:"queued"`
	Rejected int `json:"rejected"`
	Successful int `json:"successful"`
	Timedout int `json:"timedout"`
}

type ExporterConfig struct {
	Prefix string
	BrowserlessHost string
	BrowserlessPort int
	ExporterPort int
}

var (
	lastMetricsDate int64
	lastPrometheusMetricsReport string
	config ExporterConfig
)

func prefix() string {
	if len(config.Prefix) > 0 {
		return config.Prefix + "_"
	}
	return ""
}

func buildPrometheusMetricsReport(metricsObject MetricsObject) string {
	report := ""
	report += fmt.Sprintf("%sbrowserless_cpu %4.2f\n", prefix(), metricsObject.Cpu)
	report += fmt.Sprintf("%sbrowserless_memory %4.2f\n", prefix(), metricsObject.Memory)
	report += fmt.Sprintf("%sbrowserless_successful %d\n", prefix(), metricsObject.Successful)
	report += fmt.Sprintf("%sbrowserless_queued %d\n", prefix(), metricsObject.Queued)
	report += fmt.Sprintf("%sbrowserless_rejected %d\n", prefix(), metricsObject.Rejected)
	report += fmt.Sprintf("%sbrowserless_timedout %d\n", prefix(), metricsObject.Timedout)
	report += fmt.Sprintf("%sbrowserless_error %d\n", prefix(), metricsObject.Error)
	return report
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(fmt.Sprintf("http://%s:%d/metrics", config.BrowserlessHost, config.BrowserlessPort))

	if err != nil {
		log.Fatal(err)
		w.Write(make([]byte, 0))
		return
	}

	var buffer bytes.Buffer
	io.Copy(&buffer, res.Body)
	
	metricsObjects := make([]MetricsObject, 0)
	json.Unmarshal(buffer.Bytes(), &metricsObjects)

	lastMetricsObject := metricsObjects[len(metricsObjects) - 1]

	if lastMetricsObject.Date != lastMetricsDate {
		lastPrometheusMetricsReport = buildPrometheusMetricsReport(lastMetricsObject)
	}

	w.Write([]byte(lastPrometheusMetricsReport))
}

func main() {
	lastMetricsDate = 0
	lastPrometheusMetricsReport = ""

	config.Prefix = *flag.String("prefix", "", "Prefix for metrics names")
	config.BrowserlessHost = *flag.String("browserless.host", "127.0.0.1", "Browserless host")
	config.BrowserlessPort = *flag.Int("browserless.port", 3000, "Browserless port")
	config.ExporterPort = *flag.Int("exporter.port", 3002, "Exporter port")

	http.HandleFunc("/metrics", handleMetrics)

	exporterAddress := fmt.Sprintf(":%d", config.ExporterPort)

	http.ListenAndServe(exporterAddress, nil)
}

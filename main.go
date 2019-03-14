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
	Prefix *string
	BrowserlessHost *string
	BrowserlessPort *int
	ExporterHost *string
	ExporterPort *int
}

var (
	lastMetricsDate int64
	lastPrometheusMetricsReport string
	config ExporterConfig
	browserlessAddress string
)

func prefix() string {
	if len(*config.Prefix) > 0 {
		return *config.Prefix + "_"
	}
	return ""
}

func buildPrometheusMetricsReport(metricsObject MetricsObject, count int) string {
	report := ""
	report += fmt.Sprintf("%sbrowserless_cpu_sum %4.2f\n", prefix(), metricsObject.Cpu)
	report += fmt.Sprintf("%sbrowserless_cpu_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_memory_sum %4.2f\n", prefix(), metricsObject.Memory)
	report += fmt.Sprintf("%sbrowserless_memory_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_successful_sum %d\n", prefix(), metricsObject.Successful)
	report += fmt.Sprintf("%sbrowserless_successful_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_queued_sum %d\n", prefix(), metricsObject.Queued)
	report += fmt.Sprintf("%sbrowserless_queued_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_rejected_sum %d\n", prefix(), metricsObject.Rejected)
	report += fmt.Sprintf("%sbrowserless_rejected_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_timedout_sum %d\n", prefix(), metricsObject.Timedout)
	report += fmt.Sprintf("%sbrowserless_timedout_count %d\n", prefix(), count)
	report += fmt.Sprintf("%sbrowserless_error_sum %d\n", prefix(), metricsObject.Error)
	report += fmt.Sprintf("%sbrowserless_error_count %d\n", prefix(), count)
	return report
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(fmt.Sprintf("http://%s/metrics", browserlessAddress))

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
		var sumMetricsObject MetricsObject

		for _, metricsObject := range metricsObjects {
			sumMetricsObject.Cpu += metricsObject.Cpu
			sumMetricsObject.Error += metricsObject.Error
			sumMetricsObject.Memory += metricsObject.Memory
			sumMetricsObject.Queued += metricsObject.Queued
			sumMetricsObject.Rejected += metricsObject.Rejected
			sumMetricsObject.Successful += metricsObject.Successful
			sumMetricsObject.Timedout += metricsObject.Timedout
		}
		
		lastPrometheusMetricsReport = buildPrometheusMetricsReport(sumMetricsObject, len(metricsObjects))
	}

	w.Write([]byte(lastPrometheusMetricsReport))
}

func getBrowserlessAddress(config ExporterConfig) string {
	host := *config.BrowserlessHost
	port := *config.BrowserlessPort
	return fmt.Sprintf("%s:%d", host, port)
}

func getExporterAddress(config ExporterConfig) string {
	host := *config.ExporterHost
	port := *config.ExporterPort
	return fmt.Sprintf("%s:%d", host, port)
}

func main() {
	lastMetricsDate = 0
	lastPrometheusMetricsReport = ""

	config.Prefix = flag.String("prefix", "", "Prefix for metrics names")
	config.BrowserlessHost = flag.String("browserless.host", "localhost", "Browserless host")
	config.BrowserlessPort = flag.Int("browserless.port", 3000, "Browserless port")
	config.ExporterHost = flag.String("exporter.host", "localhost", "Exporter host")
	config.ExporterPort = flag.Int("exporter.port", 3002, "Exporter port")

	flag.Parse()

	browserlessAddress = getBrowserlessAddress(config)
	log.Printf("Browserless address %s\n", browserlessAddress)

	exporterAddress := getExporterAddress(config)
	log.Printf("Starting exporter on %s\n", exporterAddress)

	http.HandleFunc("/metrics", handleMetrics)

	http.ListenAndServe(exporterAddress, nil)
}

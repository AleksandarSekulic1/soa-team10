package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Global Prometheus client
var prometheusClient *PrometheusClient

// Simple metrics storage
var (
	requestCounter   int64
	errorCounter     int64
	startTime        = time.Now()
	requestDurations []float64
	activeRequests   int64
)

// initPrometheusClient initializes the Prometheus client
func initPrometheusClient() {
	// Prometheus is accessible at prometheus:9090 from within the container network
	prometheusClient = NewPrometheusClient("http://prometheus:9090")
}

// getRealSystemMetrics gets actual system metrics from both container and host
func getRealSystemMetrics() SystemMetrics {
	// Get Container metrics (from inside the Docker container)
	containerCpuUsage := getCPUUsage()
	containerMemoryUsage := getMemoryUsage()
	containerDiskUsage := getDiskUsage()
	containerNetworkRx, containerNetworkTx := getNetworkUsage()
	
	// Get Host metrics (from the Windows PC running the containers)
	hostCpuUsage := getHostCPUUsage()
	hostMemoryUsage := getHostMemoryUsage()
	hostDiskUsage := getHostDiskUsage()
	hostNetworkRx, hostNetworkTx := getHostNetworkUsage()
	
	return SystemMetrics{
		Container: ContainerMetrics{
			CPUUsage:    containerCpuUsage,
			MemoryUsage: containerMemoryUsage,
			DiskUsage:   containerDiskUsage,
			NetworkRx:   containerNetworkRx,
			NetworkTx:   containerNetworkTx,
		},
		Host: HostMetrics{
			CPUUsage:    hostCpuUsage,
			MemoryUsage: hostMemoryUsage,
			DiskUsage:   hostDiskUsage,
			NetworkRx:   hostNetworkRx,
			NetworkTx:   hostNetworkTx,
		},
	}
}

// getCPUUsage gets real CPU usage using container-compatible methods
func getCPUUsage() float64 {
	// Try to read from /proc/stat (Linux container)
	if content, err := exec.Command("cat", "/proc/stat").Output(); err == nil {
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "cpu ") {
			fields := strings.Fields(lines[0])
			if len(fields) >= 8 {
				user, _ := strconv.ParseFloat(fields[1], 64)
				nice, _ := strconv.ParseFloat(fields[2], 64)
				system, _ := strconv.ParseFloat(fields[3], 64)
				idle, _ := strconv.ParseFloat(fields[4], 64)
				
				total := user + nice + system + idle
				if total > 0 {
					usage := ((total - idle) / total) * 100
					return math.Min(usage, 100) // Cap at 100%
				}
			}
		}
	}
	
	// Container fallback: return 0 instead of random value
	return 0
}

// getMemoryUsage gets real memory usage using container-compatible methods
func getMemoryUsage() float64 {
	// Try to read from /proc/meminfo (Linux container)
	if content, err := exec.Command("cat", "/proc/meminfo").Output(); err == nil {
		lines := strings.Split(string(content), "\n")
		var total, available float64
		
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					total, _ = strconv.ParseFloat(fields[1], 64)
				}
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					available, _ = strconv.ParseFloat(fields[1], 64)
				}
			}
		}
		
		if total > 0 && available >= 0 {
			used := total - available
			return (used / total) * 100
		}
	}
	
	// Container fallback: return 0 instead of random value
	return 0
}

// getDiskUsage gets real disk usage using container-compatible methods
func getDiskUsage() float64 {
	// Try to get disk usage from df command (Linux container)
	if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				usageStr := strings.TrimSuffix(fields[4], "%")
				if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
					return usage
				}
			}
		}
	}
	
	// Container fallback: return 0 instead of random value
	return 0
}

// getNetworkUsage gets network statistics using container-compatible methods
func getNetworkUsage() (int64, int64) {
	// Try to read network stats from /proc/net/dev (Linux container)
	if content, err := exec.Command("cat", "/proc/net/dev").Output(); err == nil {
		lines := strings.Split(string(content), "\n")
		var totalRx, totalTx int64
		
		for _, line := range lines {
			if strings.Contains(line, ":") && !strings.Contains(line, "lo:") {
				fields := strings.Fields(line)
				if len(fields) >= 10 {
					if rxBytes, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						totalRx += rxBytes
					}
					if txBytes, err := strconv.ParseInt(fields[9], 10, 64); err == nil {
						totalTx += txBytes
					}
				}
			}
		}
		
		if totalRx > 0 || totalTx > 0 {
			return totalRx, totalTx
		}
	}
	
	// Container fallback: return 0 instead of random values
	return 0, 0
}

// getNetworkUsageFallback provides fallback network statistics (kept for compatibility)
func getNetworkUsageFallback() (int64, int64) {
	// Return 0 instead of random values
	return 0, 0
}

// HOST METRICS COLLECTION (via Prometheus and Node Exporter)
// getHostCPUUsage gets real CPU usage from Node Exporter via Prometheus
func getHostCPUUsage() float64 {
	if prometheusClient == nil {
		return 0
	}
	return prometheusClient.GetHostCPUUsage()
}

// getHostMemoryUsage gets real memory usage from Node Exporter via Prometheus
func getHostMemoryUsage() float64 {
	if prometheusClient == nil {
		return 0
	}
	return prometheusClient.GetHostMemoryUsage()
}

// getHostDiskUsage gets disk usage from Node Exporter via Prometheus
func getHostDiskUsage() float64 {
	if prometheusClient == nil {
		return 0
	}
	return prometheusClient.GetHostDiskUsage()
}

// getHostNetworkUsage gets network statistics from Node Exporter via Prometheus
func getHostNetworkUsage() (int64, int64) {
	if prometheusClient == nil {
		return 0, 0
	}
	return prometheusClient.GetHostNetworkUsage()
}

// Metrics structure
type Metrics struct {
	Service          string    `json:"service"`
	Uptime           string    `json:"uptime"`
	TotalRequests    int64     `json:"total_requests"`
	TotalErrors      int64     `json:"total_errors"`
	ActiveRequests   int64     `json:"active_requests"`
	AverageResponse  float64   `json:"average_response_time_ms"`
	MemoryUsage      string    `json:"memory_usage_mb"`
	Goroutines       int       `json:"goroutines"`
	Timestamp        time.Time `json:"timestamp"`
}

type ContainerMetrics struct {
	CPUUsage    float64 `json:"cpu_usage_percent"`
	MemoryUsage float64 `json:"memory_usage_percent"`
	DiskUsage   float64 `json:"disk_usage_percent"`
	NetworkRx   int64   `json:"network_rx_bytes"`
	NetworkTx   int64   `json:"network_tx_bytes"`
}

type HostMetrics struct {
	CPUUsage    float64 `json:"cpu_usage_percent"`
	MemoryUsage float64 `json:"memory_usage_percent"`
	DiskUsage   float64 `json:"disk_usage_percent"`
	NetworkRx   int64   `json:"network_rx_bytes"`
	NetworkTx   int64   `json:"network_tx_bytes"`
}

type SystemMetrics struct {
	Container ContainerMetrics `json:"container"`
	Host      HostMetrics      `json:"host"`
}

type MonitoringDashboard struct {
	Gateway Metrics       `json:"gateway"`
	System  SystemMetrics `json:"system"`
	Status  string        `json:"status"`
}

// Prometheus format metrics endpoint
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	uptime := time.Since(startTime).Seconds()
	var avgDuration float64
	if len(requestDurations) > 0 {
		sum := 0.0
		for _, d := range requestDurations {
			sum += d
		}
		avgDuration = sum / float64(len(requestDurations))
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memMB := float64(m.Alloc) / 1024 / 1024

	fmt.Fprintf(w, "# HELP gateway_uptime_seconds Total uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE gateway_uptime_seconds counter\n")
	fmt.Fprintf(w, "gateway_uptime_seconds %.2f\n", uptime)

	fmt.Fprintf(w, "# HELP gateway_requests_total Total number of requests\n")
	fmt.Fprintf(w, "# TYPE gateway_requests_total counter\n")
	fmt.Fprintf(w, "gateway_requests_total %d\n", atomic.LoadInt64(&requestCounter))

	fmt.Fprintf(w, "# HELP gateway_errors_total Total number of errors\n")
	fmt.Fprintf(w, "# TYPE gateway_errors_total counter\n")
	fmt.Fprintf(w, "gateway_errors_total %d\n", atomic.LoadInt64(&errorCounter))

	fmt.Fprintf(w, "# HELP gateway_active_requests Currently active requests\n")
	fmt.Fprintf(w, "# TYPE gateway_active_requests gauge\n")
	fmt.Fprintf(w, "gateway_active_requests %d\n", atomic.LoadInt64(&activeRequests))

	fmt.Fprintf(w, "# HELP gateway_memory_usage_mb Memory usage in MB\n")
	fmt.Fprintf(w, "# TYPE gateway_memory_usage_mb gauge\n")
	fmt.Fprintf(w, "gateway_memory_usage_mb %.2f\n", memMB)

	fmt.Fprintf(w, "# HELP gateway_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE gateway_goroutines gauge\n")
	fmt.Fprintf(w, "gateway_goroutines %d\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP gateway_avg_response_time_ms Average response time in milliseconds\n")
	fmt.Fprintf(w, "# TYPE gateway_avg_response_time_ms gauge\n")
	fmt.Fprintf(w, "gateway_avg_response_time_ms %.2f\n", avgDuration)
}

// JSON monitoring endpoint
func monitoringHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	
	var avgDuration float64
	if len(requestDurations) > 0 {
		sum := 0.0
		for _, d := range requestDurations {
			sum += d
		}
		avgDuration = sum / float64(len(requestDurations))
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memMB := float64(m.Alloc) / 1024 / 1024

	// Get real system metrics instead of random simulation
	system := getRealSystemMetrics()

	gateway := Metrics{
		Service:         "api-gateway",
		Uptime:          uptime.String(),
		TotalRequests:   atomic.LoadInt64(&requestCounter),
		TotalErrors:     atomic.LoadInt64(&errorCounter),
		ActiveRequests:  atomic.LoadInt64(&activeRequests),
		AverageResponse: avgDuration,
		MemoryUsage:     fmt.Sprintf("%.2f MB", memMB),
		Goroutines:      runtime.NumGoroutine(),
		Timestamp:       time.Now(),
	}

	dashboard := MonitoringDashboard{
		Gateway: gateway,
		System:  system,
		Status:  "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(dashboard)
}

// Host metrics handlers
func hostCPUHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	cpuUsage := getHostCPUUsage()
	response := map[string]interface{}{
		"cpu_usage_percent": cpuUsage,
		"timestamp":         time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func hostMemoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	memoryUsage := getHostMemoryUsage()
	response := map[string]interface{}{
		"memory_usage_percent": memoryUsage,
		"timestamp":            time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func hostDiskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	diskUsage := getHostDiskUsage()
	
	var totalSpace, freeSpace int64
	if prometheusClient != nil {
		totalSpace = prometheusClient.GetHostDiskTotalSpace()
		freeSpace = prometheusClient.GetHostDiskFreeSpace()
	}
	
	response := map[string]interface{}{
		"disk_usage_percent": diskUsage,
		"total_space_bytes":  totalSpace,
		"free_space_bytes":   freeSpace,
		"used_space_bytes":   totalSpace - freeSpace,
		"timestamp":          time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func hostNetworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	rxBytes, txBytes := getHostNetworkUsage()
	response := map[string]interface{}{
		"network_rx_bytes_per_sec": rxBytes,
		"network_tx_bytes_per_sec": txBytes,
		"timestamp":                time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func hostAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	cpuUsage := getHostCPUUsage()
	memoryUsage := getHostMemoryUsage()
	diskUsage := getHostDiskUsage()
	rxBytes, txBytes := getHostNetworkUsage()
	
	var totalSpace, freeSpace int64
	if prometheusClient != nil {
		totalSpace = prometheusClient.GetHostDiskTotalSpace()
		freeSpace = prometheusClient.GetHostDiskFreeSpace()
	}
	
	response := map[string]interface{}{
		"cpu": map[string]interface{}{
			"usage_percent": cpuUsage,
		},
		"memory": map[string]interface{}{
			"usage_percent": memoryUsage,
		},
		"disk": map[string]interface{}{
			"usage_percent":   diskUsage,
			"total_space_gb":  float64(totalSpace) / 1024 / 1024 / 1024,
			"free_space_gb":   float64(freeSpace) / 1024 / 1024 / 1024,
			"used_space_gb":   float64(totalSpace-freeSpace) / 1024 / 1024 / 1024,
			"total_space_bytes": totalSpace,
			"free_space_bytes":  freeSpace,
			"used_space_bytes":  totalSpace - freeSpace,
		},
		"network": map[string]interface{}{
			"rx_bytes_per_sec": rxBytes,
			"tx_bytes_per_sec": txBytes,
			"rx_mbps":          float64(rxBytes) / 1024 / 1024 * 8,
			"tx_mbps":          float64(txBytes) / 1024 / 1024 * 8,
		},
		"timestamp": time.Now(),
		"source":    "prometheus-node-exporter",
	}
	json.NewEncoder(w).Encode(response)
}
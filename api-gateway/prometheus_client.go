package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PrometheusClient handles communication with Prometheus API
type PrometheusClient struct {
	baseURL string
	client  *http.Client
}

// PrometheusResponse represents the structure of Prometheus API response
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// NewPrometheusClient creates a new Prometheus client
func NewPrometheusClient(prometheusURL string) *PrometheusClient {
	return &PrometheusClient{
		baseURL: prometheusURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// queryPrometheus executes a PromQL query and returns the result
func (pc *PrometheusClient) queryPrometheus(query string) (*PrometheusResponse, error) {
	// URL encode the query
	encodedQuery := url.QueryEscape(query)
	reqURL := fmt.Sprintf("%s/api/v1/query?query=%s", pc.baseURL, encodedQuery)
	
	resp, err := pc.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Prometheus returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	var promResp PrometheusResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}
	
	if promResp.Status != "success" {
		return nil, fmt.Errorf("Prometheus query failed with status: %s", promResp.Status)
	}
	
	return &promResp, nil
}

// getFirstValue extracts the first numeric value from Prometheus response
func (pc *PrometheusClient) getFirstValue(response *PrometheusResponse) (float64, error) {
	if len(response.Data.Result) == 0 {
		return 0, fmt.Errorf("no data in response")
	}
	
	if len(response.Data.Result[0].Value) < 2 {
		return 0, fmt.Errorf("invalid value format")
	}
	
	valueStr, ok := response.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("value is not a string")
	}
	
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value as float: %v", err)
	}
	
	return value, nil
}

// getSumValue sums all values from multiple results
func (pc *PrometheusClient) getSumValue(response *PrometheusResponse) (float64, error) {
	if len(response.Data.Result) == 0 {
		return 0, nil // Return 0 if no data
	}
	
	var sum float64
	for _, result := range response.Data.Result {
		if len(result.Value) < 2 {
			continue
		}
		
		valueStr, ok := result.Value[1].(string)
		if !ok {
			continue
		}
		
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}
		
		sum += value
	}
	
	return sum, nil
}

// GetHostCPUUsage gets CPU usage from WMI Exporter via Prometheus
func (pc *PrometheusClient) GetHostCPUUsage() float64 {
    // Query for CPU usage percentage from Windows WMI Exporter
    // This calculates: 100 - (average idle percentage across all cores)
    query := `100 - (avg by () (rate(windows_cpu_time_total{mode="idle"}[1m])) * 100)`
    
    response, err := pc.queryPrometheus(query)
    if err != nil {
        fmt.Printf("Error querying CPU usage: %v\n", err)
        return 0
    }
    
    value, err := pc.getFirstValue(response)
    if err != nil {
        fmt.Printf("Error extracting CPU usage value: %v\n", err)
        return 0
    }
    
    // Ensure the value is between 0 and 100
    if value < 0 {
        value = 0
    } else if value > 100 {
        value = 100
    }
    
    return value
}

// GetHostMemoryUsage gets memory usage from WMI Exporter via Prometheus
func (pc *PrometheusClient) GetHostMemoryUsage() float64 {
	// Query for memory usage percentage using Windows metrics
	// Fixed query: (total - available) / total * 100
	query := `(windows_memory_physical_total_bytes - windows_memory_available_bytes) / windows_memory_physical_total_bytes * 100`
	
	response, err := pc.queryPrometheus(query)
	if err != nil {
		fmt.Printf("Error querying memory usage: %v\n", err)
		return 0
	}
	
	value, err := pc.getFirstValue(response)
	if err != nil {
		fmt.Printf("Error extracting memory usage value: %v\n", err)
		return 0
	}
	
	// Ensure the value is between 0 and 100
	if value < 0 {
		value = 0
	} else if value > 100 {
		value = 100
	}
	
	return value
}

// GetHostDiskUsage gets disk usage from WMI Exporter via Prometheus
func (pc *PrometheusClient) GetHostDiskUsage() float64 {
	// Query for disk usage percentage using Windows logical disk metrics
	query := `max((1 - (windows_logical_disk_free_bytes{volume!~"HarddiskVolume.*"} / windows_logical_disk_size_bytes{volume!~"HarddiskVolume.*"})) * 100)`
	
	response, err := pc.queryPrometheus(query)
	if err != nil {
		fmt.Printf("Error querying disk usage: %v\n", err)
		return 0
	}
	
	value, err := pc.getFirstValue(response)
	if err != nil {
		fmt.Printf("Error extracting disk usage value: %v\n", err)
		return 0
	}
	
	// Ensure the value is between 0 and 100
	if value < 0 {
		value = 0
	} else if value > 100 {
		value = 100
	}
	
	return value
}

// GetHostNetworkUsage gets network usage from WMI Exporter via Prometheus
func (pc *PrometheusClient) GetHostNetworkUsage() (int64, int64) {
	// Query for total network bytes received from Windows metrics
	rxQuery := `sum(rate(windows_net_bytes_received_total[5m]) * 8)`
	
	rxResponse, err := pc.queryPrometheus(rxQuery)
	if err != nil {
		fmt.Printf("Error querying network RX: %v\n", err)
		return 0, 0
	}
	
	rxValue, err := pc.getFirstValue(rxResponse)
	if err != nil {
		fmt.Printf("Error extracting network RX value: %v\n", err)
		rxValue = 0
	}
	
	// Query for total network bytes transmitted from Windows metrics
	txQuery := `sum(rate(windows_net_bytes_sent_total[5m]) * 8)`
	
	txResponse, err := pc.queryPrometheus(txQuery)
	if err != nil {
		fmt.Printf("Error querying network TX: %v\n", err)
		return int64(rxValue), 0
	}
	
	txValue, err := pc.getFirstValue(txResponse)
	if err != nil {
		fmt.Printf("Error extracting network TX value: %v\n", err)
		txValue = 0
	}
	
	return int64(rxValue), int64(txValue)
}

// GetHostDiskTotalSpace gets total disk space in bytes from WMI Exporter
func (pc *PrometheusClient) GetHostDiskTotalSpace() int64 {
	query := `max(windows_logical_disk_size_bytes{volume!~"HarddiskVolume.*"})`
	
	response, err := pc.queryPrometheus(query)
	if err != nil {
		fmt.Printf("Error querying disk total space: %v\n", err)
		return 0
	}
	
	value, err := pc.getFirstValue(response)
	if err != nil {
		fmt.Printf("Error extracting disk total space value: %v\n", err)
		return 0
	}
	
	return int64(value)
}

// GetHostDiskFreeSpace gets available disk space in bytes from WMI Exporter
func (pc *PrometheusClient) GetHostDiskFreeSpace() int64 {
	query := `max(windows_logical_disk_free_bytes{volume!~"HarddiskVolume.*"})`
	
	response, err := pc.queryPrometheus(query)
	if err != nil {
		fmt.Printf("Error querying disk free space: %v\n", err)
		return 0
	}
	
	value, err := pc.getFirstValue(response)
	if err != nil {
		fmt.Printf("Error extracting disk free space value: %v\n", err)
		return 0
	}
	
	return int64(value)
}
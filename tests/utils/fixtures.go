package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestAPData provides sample AP data for testing (generic format)
func TestAPData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "AP-Test-001",
			"location":    "Building A Floor 1",
			"mac_address": "aa:bb:cc:dd:ee:01",
			"model":       "C9130AXI-A",
			"state":       "Registered",
			"ap_mode":     "Local",
			"ip_address":  "192.168.1.101",
		},
		{
			"name":        "AP-Test-002",
			"location":    "Building B Floor 2",
			"mac_address": "aa:bb:cc:dd:ee:02",
			"model":       "C9120AXI-A",
			"state":       "Downloading",
			"ap_mode":     "FlexConnect",
			"ip_address":  "192.168.1.102",
		},
		{
			"name":        "AP-Test-003",
			"location":    "Building C Floor 3",
			"mac_address": "aa:bb:cc:dd:ee:03",
			"model":       "C9115AXI-A",
			"state":       "Registered",
			"ap_mode":     "Local",
			"ip_address":  "192.168.1.103",
		},
	}
}

// TestAPTagData provides sample AP tag data for testing
func TestAPTagData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"policy_tag_name": "policy-tag-building-a",
			"site_tag_name":   "site-tag-building-a",
			"rf_tag_name":     "rf-tag-default",
		},
		{
			"policy_tag_name": "policy-tag-building-b",
			"site_tag_name":   "site-tag-building-b",
			"rf_tag_name":     "rf-tag-high-density",
		},
		{
			"policy_tag_name": "policy-tag-guest",
			"site_tag_name":   "site-tag-guest",
			"rf_tag_name":     "rf-tag-default",
		},
	}
}

// TestClientData provides sample client data for testing
func TestClientData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"mac_address":           "11:22:33:44:55:01",
			"username":              "user001",
			"ap_name":               "AP-Test-001",
			"wlan_id":               1,
			"current_switch_port":   "GigabitEthernet1/0/1",
			"vlan_id":               100,
			"session_timeout":       3600,
			"policy_type":           "WPA2",
			"encryption_cipher":     "AES",
			"authentication_method": "802.1X",
		},
		{
			"mac_address":           "11:22:33:44:55:02",
			"username":              "guest001",
			"ap_name":               "AP-Test-002",
			"wlan_id":               2,
			"current_switch_port":   "GigabitEthernet1/0/2",
			"vlan_id":               200,
			"session_timeout":       1800,
			"policy_type":           "Open",
			"encryption_cipher":     "None",
			"authentication_method": "None",
		},
	}
}

// TestWLANData provides sample WLAN data for testing
func TestWLANData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"wlan_id":             1,
			"profile_name":        "Corporate-WLAN",
			"ssid":                "Corporate",
			"admin_status":        true,
			"broadcast_ssid":      true,
			"security_type":       "WPA2",
			"authentication_list": []string{"802.1X"},
			"vlan_id":             100,
			"session_timeout":     3600,
		},
		{
			"wlan_id":             2,
			"profile_name":        "Guest-WLAN",
			"ssid":                "Guest",
			"admin_status":        true,
			"broadcast_ssid":      true,
			"security_type":       "Open",
			"authentication_list": []string{},
			"vlan_id":             200,
			"session_timeout":     1800,
		},
		{
			"wlan_id":             3,
			"profile_name":        "IOT-WLAN",
			"ssid":                "IoT-Devices",
			"admin_status":        false,
			"broadcast_ssid":      false,
			"security_type":       "WPA3",
			"authentication_list": []string{"PSK"},
			"vlan_id":             300,
			"session_timeout":     7200,
		},
	}
}

// TestOverviewData provides sample overview data for testing
func TestOverviewData() map[string]interface{} {
	return map[string]interface{}{
		"hostname":           "WLC-Test-001",
		"version":            "17.12.01",
		"uptime":             "123 days, 4 hours, 30 minutes",
		"total_aps":          25,
		"registered_aps":     23,
		"downloading_aps":    2,
		"total_clients":      156,
		"wired_clients":      89,
		"wireless_clients":   67,
		"total_wlans":        5,
		"enabled_wlans":      4,
		"cpu_utilization":    "15%",
		"memory_utilization": "45%",
		"disk_utilization":   "32%",
	}
}

// JSONFixture represents a generic JSON test fixture
type JSONFixture struct {
	Data interface{}
	JSON string
}

// CreateJSONFixture creates a JSON fixture from any data structure
func CreateJSONFixture(t *testing.T, data interface{}) *JSONFixture {
	t.Helper()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data to JSON: %v", err)
	}

	return &JSONFixture{
		Data: data,
		JSON: string(jsonData),
	}
}

// ExpectedTableHeaders provides common table headers for validation
func ExpectedTableHeaders() map[string][]string {
	return map[string][]string{
		"ap": {
			"NAME",
			"LOCATION",
			"MAC ADDRESS",
			"MODEL",
			"STATE",
			"AP MODE",
			"IP ADDRESS",
		},
		"ap-tag": {
			"POLICY TAG NAME",
			"SITE TAG NAME",
			"RF TAG NAME",
		},
		"client": {
			"MAC ADDRESS",
			"USERNAME",
			"AP NAME",
			"WLAN ID",
			"CURRENT SWITCH PORT",
			"VLAN ID",
			"SESSION TIMEOUT",
			"POLICY TYPE",
			"ENCRYPTION CIPHER",
			"AUTHENTICATION METHOD",
		},
		"wlan": {
			"WLAN ID",
			"PROFILE NAME",
			"SSID",
			"ADMIN STATUS",
			"BROADCAST SSID",
			"SECURITY TYPE",
			"AUTHENTICATION LIST",
			"VLAN ID",
			"SESSION TIMEOUT",
		},
		"overview": {
			"PROPERTY",
			"VALUE",
		},
	}
}

// MockResponse represents a mock API response
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// CreateMockResponse creates a mock HTTP response
func CreateMockResponse(statusCode int, data interface{}) *MockResponse {
	jsonData, _ := json.Marshal(data)
	return &MockResponse{
		StatusCode: statusCode,
		Body:       string(jsonData),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// ErrorResponse represents various error scenarios
type ErrorResponse struct {
	Type    string
	Code    int
	Message string
}

// CommonErrorResponses provides standard error responses for testing
func CommonErrorResponses() map[string]*ErrorResponse {
	return map[string]*ErrorResponse{
		"unauthorized": {
			Type:    "Authentication Error",
			Code:    401,
			Message: "Invalid credentials",
		},
		"forbidden": {
			Type:    "Authorization Error",
			Code:    403,
			Message: "Access denied",
		},
		"not_found": {
			Type:    "Resource Not Found",
			Code:    404,
			Message: "Resource not found",
		},
		"timeout": {
			Type:    "Timeout Error",
			Code:    408,
			Message: "Request timeout",
		},
		"server_error": {
			Type:    "Internal Server Error",
			Code:    500,
			Message: "Internal server error occurred",
		},
		"bad_gateway": {
			Type:    "Gateway Error",
			Code:    502,
			Message: "Bad gateway",
		},
		"service_unavailable": {
			Type:    "Service Unavailable",
			Code:    503,
			Message: "Service temporarily unavailable",
		},
	}
}

// TestControllerEndpoints provides test controller endpoints
func TestControllerEndpoints() map[string]string {
	return map[string]string{
		"primary":   "https://controller1.example.com",
		"secondary": "https://controller2.example.com",
		"local":     "https://localhost:8443",
		"invalid":   "not-a-valid-url",
	}
}

// ValidateJSONOutput validates that output is valid JSON and contains expected fields
func ValidateJSONOutput(t *testing.T, output string, expectedFields ...string) {
	t.Helper()

	var result interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Check if it's an array or object
	switch data := result.(type) {
	case []interface{}:
		if len(data) == 0 {
			return // Empty array is valid
		}
		// Check first item for expected fields
		if item, ok := data[0].(map[string]interface{}); ok {
			validateFields(t, item, expectedFields...)
		}
	case map[string]interface{}:
		validateFields(t, data, expectedFields...)
	default:
		t.Errorf("Expected JSON object or array, got %T", result)
	}
}

// validateFields checks if all expected fields exist in the data
func validateFields(t *testing.T, data map[string]interface{}, expectedFields ...string) {
	t.Helper()

	for _, field := range expectedFields {
		if _, exists := data[field]; !exists {
			t.Errorf("Expected field %q not found in JSON output", field)
		}
	}
}

// FormatTestName creates a standardized test name
func FormatTestName(category, subtest, scenario string) string {
	return fmt.Sprintf("Test%s_%s_%s", category, subtest, scenario)
}

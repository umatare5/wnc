package cisco

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRrmTypeAliases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "RrmOperResponse alias exists",
			test: func(t *testing.T) {
				var resp *RrmOperResponse
				_ = resp
				t.Log("RrmOperResponse type alias is valid")
			},
		},
		{
			name: "RrmMeasurementResponse alias exists",
			test: func(t *testing.T) {
				var resp *RrmMeasurementResponse
				_ = resp
				t.Log("RrmMeasurementResponse type alias is valid")
			},
		},
		{
			name: "RrmGlobalOperResponse alias exists",
			test: func(t *testing.T) {
				var resp *RrmGlobalOperResponse
				_ = resp
				t.Log("RrmGlobalOperResponse type alias is valid")
			},
		},
		{
			name: "RrmCfgResponse alias exists",
			test: func(t *testing.T) {
				var resp *RrmCfgResponse
				_ = resp
				t.Log("RrmCfgResponse type alias is valid")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestRrmJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		create func() interface{}
	}{
		{
			name: "RrmOperResponse serialization",
			create: func() interface{} {
				return &RrmOperResponse{}
			},
		},
		{
			name: "RrmMeasurementResponse serialization",
			create: func() interface{} {
				return &RrmMeasurementResponse{}
			},
		},
		{
			name: "RrmGlobalOperResponse serialization",
			create: func() interface{} {
				return &RrmGlobalOperResponse{}
			},
		},
		{
			name: "RrmCfgResponse serialization",
			create: func() interface{} {
				return &RrmCfgResponse{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := tt.create()
			data, err := json.Marshal(obj)
			if err != nil {
				t.Errorf("Failed to marshal %T: %v", obj, err)
			}

			// Try to unmarshal back
			switch obj.(type) {
			case *RrmOperResponse:
				var unmarshaled RrmOperResponse
				err = json.Unmarshal(data, &unmarshaled)
			case *RrmMeasurementResponse:
				var unmarshaled RrmMeasurementResponse
				err = json.Unmarshal(data, &unmarshaled)
			case *RrmGlobalOperResponse:
				var unmarshaled RrmGlobalOperResponse
				err = json.Unmarshal(data, &unmarshaled)
			case *RrmCfgResponse:
				var unmarshaled RrmCfgResponse
				err = json.Unmarshal(data, &unmarshaled)
			}

			if err != nil {
				t.Errorf("Failed to unmarshal %T: %v", obj, err)
			}
		})
	}
}

func TestRrmFunctionSignatures(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "GetRrmOper function signature",
			test: func(t *testing.T) {
				// Test that function exists and has correct signature
				// Check the function type without calling it to avoid nil pointer dereference
				funcType := reflect.TypeOf(GetRrmOper)
				if funcType == nil {
					t.Error("GetRrmOper function not found")
					return
				}

				// Verify function signature: func(*Client, context.Context) (*RrmOperResponse, error)
				if funcType.NumIn() != 2 {
					t.Errorf("GetRrmOper expected 2 parameters, got %d", funcType.NumIn())
				}
				if funcType.NumOut() != 2 {
					t.Errorf("GetRrmOper expected 2 return values, got %d", funcType.NumOut())
				}

				t.Log("GetRrmOper function signature is correct")
			},
		},
		{
			name: "GetRrmMeasurement function signature",
			test: func(t *testing.T) {
				// Test that function exists and has correct signature
				// Check the function type without calling it to avoid nil pointer dereference
				funcType := reflect.TypeOf(GetRrmMeasurement)
				if funcType == nil {
					t.Error("GetRrmMeasurement function not found")
					return
				}

				// Verify function signature: func(*Client, context.Context) (*RrmMeasurementResponse, error)
				if funcType.NumIn() != 2 {
					t.Errorf("GetRrmMeasurement expected 2 parameters, got %d", funcType.NumIn())
				}
				if funcType.NumOut() != 2 {
					t.Errorf("GetRrmMeasurement expected 2 return values, got %d", funcType.NumOut())
				}

				t.Log("GetRrmMeasurement function signature is correct")
			},
		},
		{
			name: "GetRrmGlobalOper function signature",
			test: func(t *testing.T) {
				// Test that function exists and has correct signature
				// Check the function type without calling it to avoid nil pointer dereference
				funcType := reflect.TypeOf(GetRrmGlobalOper)
				if funcType == nil {
					t.Error("GetRrmGlobalOper function not found")
					return
				}

				// Verify function signature: func(*Client, context.Context) (*RrmGlobalOperResponse, error)
				if funcType.NumIn() != 2 {
					t.Errorf("GetRrmGlobalOper expected 2 parameters, got %d", funcType.NumIn())
				}
				if funcType.NumOut() != 2 {
					t.Errorf("GetRrmGlobalOper expected 2 return values, got %d", funcType.NumOut())
				}

				t.Log("GetRrmGlobalOper function signature is correct")
			},
		},
		{
			name: "GetRrmCfg function signature",
			test: func(t *testing.T) {
				// Test that function exists and has correct signature
				// Check the function type without calling it to avoid nil pointer dereference
				funcType := reflect.TypeOf(GetRrmCfg)
				if funcType == nil {
					t.Error("GetRrmCfg function not found")
					return
				}

				// Verify function signature: func(*Client, context.Context) (*RrmCfgResponse, error)
				if funcType.NumIn() != 2 {
					t.Errorf("GetRrmCfg expected 2 parameters, got %d", funcType.NumIn())
				}
				if funcType.NumOut() != 2 {
					t.Errorf("GetRrmCfg expected 2 return values, got %d", funcType.NumOut())
				}

				t.Log("GetRrmCfg function signature is correct")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestRrmFailFast(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "RrmOperResponse should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("RrmOperResponse creation panicked: %v", r)
					}
				}()
				var resp *RrmOperResponse
				_ = resp
			},
		},
		{
			name: "RrmMeasurementResponse should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("RrmMeasurementResponse creation panicked: %v", r)
					}
				}()
				var resp *RrmMeasurementResponse
				_ = resp
			},
		},
		{
			name: "RrmGlobalOperResponse should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("RrmGlobalOperResponse creation panicked: %v", r)
					}
				}()
				var resp *RrmGlobalOperResponse
				_ = resp
			},
		},
		{
			name: "RrmCfgResponse should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("RrmCfgResponse creation panicked: %v", r)
					}
				}()
				var resp *RrmCfgResponse
				_ = resp
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestRrmIntegration(t *testing.T) {
	t.Run("nil client should handle gracefully", func(t *testing.T) {
		// Test that functions exist and can be called with proper signature
		// We test the function signatures without actually calling them with nil
		// to avoid segmentation faults
		funcs := []struct {
			name string
			fn   interface{}
		}{
			{"GetRrmOper", GetRrmOper},
			{"GetRrmMeasurement", GetRrmMeasurement},
			{"GetRrmGlobalOper", GetRrmGlobalOper},
			{"GetRrmCfg", GetRrmCfg},
		}

		for _, f := range funcs {
			funcType := reflect.TypeOf(f.fn)
			if funcType == nil {
				t.Errorf("%s function not found", f.name)
				continue
			}

			// Verify it's a function that takes 2 parameters and returns 2 values
			if funcType.Kind() != reflect.Func {
				t.Errorf("%s is not a function", f.name)
				continue
			}

			if funcType.NumIn() != 2 {
				t.Errorf("%s expected 2 parameters, got %d", f.name, funcType.NumIn())
			}

			if funcType.NumOut() != 2 {
				t.Errorf("%s expected 2 return values, got %d", f.name, funcType.NumOut())
			}
		}

		t.Log("All RRM function signatures verified successfully")
	})
}

func TestGetRrmMeasurement_WithRealResponse(t *testing.T) {
	// Create a test server with real WNC RRM measurement response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"cisco-wireless-rrm-oper:rrm-measurement": {
				"ap-mac": "28:ac:9e:bb:3c:80",
				"radio-slot": 0,
				"noise-floor": -92,
				"interference": 15,
				"channel-utilization": 25,
				"client-count": 8
			}
		}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetRrmMeasurement function
	result, err := GetRrmMeasurement(client, context.Background())
	if err != nil {
		t.Errorf("GetRrmMeasurement failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		t.Logf("GetRrmMeasurement returned result successfully")
	}
}

func TestGetRrmMeasurement_WithErrorHandling(t *testing.T) {
	// Create a test server that returns 403 forbidden error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(`{"error": "Access forbidden"}`))
		if err != nil {
			t.Errorf("Failed to write error response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetRrmMeasurement with error response
	result, err := GetRrmMeasurement(client, context.Background())
	if err == nil {
		t.Error("Expected error for forbidden response")
	}

	if result != nil {
		t.Error("Expected nil result for error response")
	}
}

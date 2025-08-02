package mock

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/umatare5/wnc/pkg/cisco"
)

// TestMockWNCClient tests mock WNC client functionality (Mock test)
func TestMockWNCClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockWNCClient(ctrl)
	ctx := context.Background()

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "mock_get_ap_oper",
			test: func(t *testing.T) {
				// Set up mock expectation
				expectedResponse := &cisco.ApOperResponse{}
				mockClient.EXPECT().
					GetApOper(ctx).
					Return(expectedResponse, nil).
					Times(1)

				// Execute the method
				response, err := mockClient.GetApOper(ctx)

				// Verify results
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if response != expectedResponse {
					t.Errorf("Expected response %v, got %v", expectedResponse, response)
				}
			},
		},
		{
			name: "mock_get_wlan_cfg",
			test: func(t *testing.T) {
				// Set up mock expectation
				expectedResponse := &cisco.WlanCfgResponse{}
				mockClient.EXPECT().
					GetWlanCfg(ctx).
					Return(expectedResponse, nil).
					Times(1)

				// Execute the method
				response, err := mockClient.GetWlanCfg(ctx)

				// Verify results
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if response != expectedResponse {
					t.Errorf("Expected response %v, got %v", expectedResponse, response)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestMockClientFactory tests mock client factory functionality (Mock test)
func TestMockClientFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := NewMockClientFactory(ctrl)

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "mock_new_client",
			test: func(t *testing.T) {
				// Set up mock expectation
				var isSecure *bool = nil
				mockClient := NewMockWNCClient(ctrl)
				mockFactory.EXPECT().
					NewClient("test-controller", "test-token", isSecure).
					Return(mockClient, nil).
					Times(1)

				// Execute the method
				client, err := mockFactory.NewClient("test-controller", "test-token", isSecure)

				// Verify results
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if client != mockClient {
					t.Errorf("Expected client %v, got %v", mockClient, client)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

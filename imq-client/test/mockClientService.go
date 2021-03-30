package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/protocol"
)

// MockClientService is a struct for mocking ClientService
type MockClientService struct {
	Mock
	client.Service
}

// GetID mocks on ClientService.GetID
func (m *MockClientService) GetID() int {
	args := m.Called()
	return args.Int(0)
}

// GetAddress mocks on ClientService.GetAddress
func (m *MockClientService) GetAddress() string {
	args := m.Called()
	return args.String(0)
}

// SendRequest mocks on ClientService.SendRequest
func (m *MockClientService) SendRequest(ctx context.Context, request *protocol.Request) ([]byte, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]byte), args.Error(1)
}

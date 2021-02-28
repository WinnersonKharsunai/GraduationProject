package test

import messagefactory "github.com/WinnersonKharsunai/GraduationProject/client/message-factory"

// MockMessageFactoryIF is a struct for mocking MessageFactoryIF
type MockMessageFactoryIF struct {
	Mock
	messagefactory.MessagefactoryIF
}

// MarshalRequestBody mocks on MessageFactoryIF.MarshalRequestBody
func (m *MockMessageFactoryIF) MarshalRequestBody(v interface{}, contentType string) ([]byte, error) {
	args := m.Called(v, contentType)
	return args.Get(0).([]byte), args.Error(1)
}

// UnmarshalRequestBody mocks on MessageFactoryIF.UnmarshalRequestBody
func (m *MockMessageFactoryIF) UnmarshalRequestBody(data []byte, v interface{}, contentType string) error {
	args := m.Called(data, v, contentType)
	return args.Error(0)
}

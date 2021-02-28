package test

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/stretchr/testify/mock"
)

// Mock type
type Mock struct {
	mock.Mock
}

// MockSetup type is concrete implementation for mocking
type MockSetup struct {
	*mock.Call
	methodName string
	mock       *Mock
	args       []interface{}
	returns    []interface{}
}

// Given sets the given method name
func (mock *Mock) Given(methodName interface{}) *MockSetup {
	name := determineMethodname(methodName)
	return &MockSetup{
		methodName: name,
		mock:       mock,
	}
}

// When sets the given arguments
func (m *MockSetup) When(args ...interface{}) *MockSetup {
	m.args = append(m.args, args)
	m.Call = m.mock.On(m.methodName, args...)
	return m
}

// Then sets the given returns values
func (m *MockSetup) Then(returns ...interface{}) *MockSetup {
	m.returns = append(m.returns, returns)
	m.Call.Return(returns...)
	return m
}

func determineMethodname(v interface{}) string {
	if v == nil {
		panic("must provide a non-nil value for method name")
	}
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		panic("must provide a non-pointer value for method name")
	}

	switch v.(type) {
	case string:
		return v.(string)
	default:
		return getFuncionName(v)
	}
}

func getFuncionName(v interface{}) string {
	fullname := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	splitName := strings.Split(fullname, ".")
	return splitName[len(splitName)-1]
}

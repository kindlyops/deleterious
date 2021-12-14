// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	mock "github.com/stretchr/testify/mock"
)

// DynamoDBAPI is an autogenerated mock type for the DynamoDBAPI type
type DynamoDBAPI struct {
	mock.Mock
}

// ListTables provides a mock function with given fields: params
func (_m *DynamoDBAPI) ListTables(params *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	ret := _m.Called(params)

	var r0 *dynamodb.ListTablesOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.ListTablesInput) *dynamodb.ListTablesOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ListTablesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.ListTablesInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	kms "github.com/aws/aws-sdk-go/service/kms"
	mock "github.com/stretchr/testify/mock"
)

// KmsAPI is an autogenerated mock type for the KmsAPI type
type KmsAPI struct {
	mock.Mock
}

// DescribeKey provides a mock function with given fields: params
func (_m *KmsAPI) DescribeKey(params *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	ret := _m.Called(params)

	var r0 *kms.DescribeKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.DescribeKeyInput) *kms.DescribeKeyOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DescribeKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DescribeKeyInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListKeys provides a mock function with given fields: params
func (_m *KmsAPI) ListKeys(params *kms.ListKeysInput) (*kms.ListKeysOutput, error) {
	ret := _m.Called(params)

	var r0 *kms.ListKeysOutput
	if rf, ok := ret.Get(0).(func(*kms.ListKeysInput) *kms.ListKeysOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListKeysOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListKeysInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

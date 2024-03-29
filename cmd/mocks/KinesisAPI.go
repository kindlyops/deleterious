// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	kinesis "github.com/aws/aws-sdk-go/service/kinesis"
	mock "github.com/stretchr/testify/mock"
)

// KinesisAPI is an autogenerated mock type for the KinesisAPI type
type KinesisAPI struct {
	mock.Mock
}

// ListStreams provides a mock function with given fields: params
func (_m *KinesisAPI) ListStreams(params *kinesis.ListStreamsInput) (*kinesis.ListStreamsOutput, error) {
	ret := _m.Called(params)

	var r0 *kinesis.ListStreamsOutput
	if rf, ok := ret.Get(0).(func(*kinesis.ListStreamsInput) *kinesis.ListStreamsOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kinesis.ListStreamsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kinesis.ListStreamsInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	cloudwatchlogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	mock "github.com/stretchr/testify/mock"
)

// LogsAPI is an autogenerated mock type for the LogsAPI type
type LogsAPI struct {
	mock.Mock
}

// DescribeLogGroups provides a mock function with given fields: params
func (_m *LogsAPI) DescribeLogGroups(params *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	ret := _m.Called(params)

	var r0 *cloudwatchlogs.DescribeLogGroupsOutput
	if rf, ok := ret.Get(0).(func(*cloudwatchlogs.DescribeLogGroupsInput) *cloudwatchlogs.DescribeLogGroupsOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*cloudwatchlogs.DescribeLogGroupsInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeLogStreams provides a mock function with given fields: params
func (_m *LogsAPI) DescribeLogStreams(params *cloudwatchlogs.DescribeLogStreamsInput) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	ret := _m.Called(params)

	var r0 *cloudwatchlogs.DescribeLogStreamsOutput
	if rf, ok := ret.Get(0).(func(*cloudwatchlogs.DescribeLogStreamsInput) *cloudwatchlogs.DescribeLogStreamsOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudwatchlogs.DescribeLogStreamsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*cloudwatchlogs.DescribeLogStreamsInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

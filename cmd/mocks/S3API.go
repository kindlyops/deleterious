// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	s3 "github.com/aws/aws-sdk-go/service/s3"
	mock "github.com/stretchr/testify/mock"
)

// S3API is an autogenerated mock type for the S3API type
type S3API struct {
	mock.Mock
}

// ListBuckets provides a mock function with given fields: params
func (_m *S3API) ListBuckets(params *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	ret := _m.Called(params)

	var r0 *s3.ListBucketsOutput
	if rf, ok := ret.Get(0).(func(*s3.ListBucketsInput) *s3.ListBucketsOutput); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.ListBucketsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*s3.ListBucketsInput) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

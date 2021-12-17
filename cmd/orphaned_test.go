// Copyright © 2021 Kindly Ops, LLC <support@kindlyops.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kindlyops/deleterious/cmd/mocks"
)

func Test_getRootedResources(t *testing.T) {
	t.Parallel()

	type args struct {
		svc  *mocks.CloudformationAPI
		kind string
	}

	mockSvc := &mocks.CloudformationAPI{}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			name: "test",
			args: args{svc: mockSvc, kind: "AWS::S3::Bucket"},
			want: map[string]bool{
				"FirstBucket": true,
			},
		},
	}

	mockSvc.On("ListStacks",
		&cloudformation.ListStacksInput{
			StackStatusFilter: getStackStates(),
		}).Return(&cloudformation.ListStacksOutput{
		StackSummaries: []*cloudformation.StackSummary{
			{
				StackName:   aws.String("FIRST_STACK"),
				StackStatus: aws.String("CREATE_COMPLETE"),
			},
			{
				StackName:   aws.String("SECOND_STACK"),
				StackStatus: aws.String("CREATE_COMPLETE"),
			},
		}}, nil).Twice()
	mockSvc.On("ListStackResources",
		&cloudformation.ListStackResourcesInput{
			StackName: aws.String("FIRST_STACK"),
		}).Return(&cloudformation.ListStackResourcesOutput{
		StackResourceSummaries: []*cloudformation.StackResourceSummary{
			{
				ResourceType:       aws.String("AWS::S3::Bucket"),
				LogicalResourceId:  aws.String("FirstBucket"),
				PhysicalResourceId: aws.String("FirstBucket"),
			},
		}}, nil).Once()
	mockSvc.On("ListStackResources",
		&cloudformation.ListStackResourcesInput{
			StackName: aws.String("SECOND_STACK"),
		}).Return(&cloudformation.ListStackResourcesOutput{
		StackResourceSummaries: []*cloudformation.StackResourceSummary{
			{
				ResourceType:       aws.String("AWS::EC2::Instance"),
				LogicalResourceId:  aws.String("OtherResource"),
				PhysicalResourceId: aws.String("OtherResource"),
			},
		}}, nil).Once()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getRootedResources(tt.args.svc, tt.args.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRootedResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

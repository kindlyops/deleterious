// Copyright © 2018 Kindly Ops, LLC <support@kindlyops.com>
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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/spf13/cobra"
)

// orphanedCmd represents the orphaned command
var orphanedCmd = &cobra.Command{
	Use:   "orphaned",
	Short: "Find orphaned resources without deleting",
	Long:  `Finds all of the specified resources that are not referenced from an active CloudFormation stack.`,
	Run:   orphaned,
}

func orphaned(cmd *cobra.Command, args []string) {

	svc := cloudformation.New(getSession())

	states := []*string{
		aws.String("CREATE_IN_PROGRESS"),
		aws.String("CREATE_FAILED"),
		aws.String("CREATE_COMPLETE"),
		aws.String("ROLLBACK_IN_PROGRESS"),
		aws.String("ROLLBACK_FAILED"),
		aws.String("ROLLBACK_COMPLETE"),
		aws.String("DELETE_IN_PROGRESS"),
		aws.String("DELETE_FAILED"),
		// aws.String("DELETE_COMPLETE"), // deleted stacks hang around for 90 days
		aws.String("UPDATE_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE"),
		aws.String("UPDATE_ROLLBACK_IN_PROGRESS"),
		aws.String("UPDATE_ROLLBACK_FAILED"),
		aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"),
		aws.String("UPDATE_ROLLBACK_COMPLETE"),
		aws.String("REVIEW_IN_PROGRESS"),
	}

	rootedResources := make(map[string]bool)
	var token *string
	for ok := true; ok; ok = (token != nil) {
		stacks, err := svc.ListStacks(&cloudformation.ListStacksInput{
			NextToken:         token,
			StackStatusFilter: states,
		})
		if err != nil {
			fmt.Printf("Error listing stacks: %v", err)
		} else {
			token = stacks.NextToken
		}

		for _, stack := range stacks.StackSummaries {
			if Debug {
				fmt.Printf("Processing %v\n", *stack.StackName)
			}
			resources, err := svc.ListStackResources(&cloudformation.ListStackResourcesInput{
				StackName: stack.StackName,
			})
			if err != nil {
				fmt.Printf("Error processing %v\n", *stack.StackName)
				continue
			}
			for _, resource := range resources.StackResourceSummaries {
				if *resource.ResourceType == Resource {
					if Debug {
						fmt.Printf("Found rooted resource %v : %v\n", *resource.PhysicalResourceId, *resource.ResourceType)
					}
					rootedResources[*resource.PhysicalResourceId] = true

				}
			}
		}
	}

	switch {
	case Resource == "AWS::DynamoDB::Table":
		processDynamoDB(rootedResources)
	case Resource == "AWS::KMS::Key":
		processKMS(rootedResources)
	case Resource == "AWS::Kinesis::Stream":
		processKinesis(rootedResources)
	default:
		fmt.Printf("%v is not yet supported", Resource)
	}
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}))
}

func processKinesis(rootedResources map[string]bool) {
	svc := kinesis.New(getSession())

	var startStream *string
	for ok := true; ok; ok = (startStream != nil) {
		streams, err := svc.ListStreams(&kinesis.ListStreamsInput{
			ExclusiveStartStreamName: startStream,
			Limit:                    aws.Int64(50),
		})
		if err != nil {
			fmt.Printf("Error listing Kinesis Streams: %v\n", err)
		} else {
			if *streams.HasMoreStreams {
				if Debug {
					fmt.Printf("Has more streams, setting startStream\n")
				}
				startStream = streams.StreamNames[len(streams.StreamNames)-1]
			} else {
				startStream = nil
			}
		}
		for _, stream := range streams.StreamNames {
			if Debug {
				fmt.Printf("Processing %v\n", *stream)
			}
			if _, ok := rootedResources[*stream]; ok {
				// this stream is owned by a cloudformation stack, skip it
				if Debug {
					fmt.Printf("Stream %v is owned by a cloudformation stack, skipping\n", *stream)
				}
			} else {
				fmt.Printf("\"%v\"\n", *stream)
			}
		}
	}
}

func processKMS(rootedResources map[string]bool) {
	svc := kms.New(getSession())

	var marker *string
	for ok := true; ok; ok = (marker != nil) {
		keys, err := svc.ListKeys(&kms.ListKeysInput{
			Marker: marker,
		})
		if err != nil {
			fmt.Printf("Error listing KMS keys: %v\n", err)
		} else {
			marker = keys.NextMarker
		}

		for _, key := range keys.Keys {
			if _, ok := rootedResources[*key.KeyId]; ok {
				// this key is owned by a cloudformation stack, skip it
			} else {
				fmt.Printf("\"%v\"\n", *key.KeyId)
			}
		}
	}
}

func processDynamoDB(rootedResources map[string]bool) {
	svc := dynamodb.New(getSession())

	var lastTable *string
	for ok := true; ok; ok = (lastTable != nil) {
		tables, err := svc.ListTables(&dynamodb.ListTablesInput{
			ExclusiveStartTableName: lastTable,
		})
		if err != nil {
			fmt.Printf("Error listing DynamoDB tables: %v\n", err)
		} else {
			lastTable = tables.LastEvaluatedTableName
		}

		//fmt.Printf("Processing %v tables, %v will be excluded\n", len(tables.TableNames), len(rootedResources))
		for _, table := range tables.TableNames {
			if _, ok := rootedResources[*table]; ok {
				// this table table belongs to a cloudformation stack, skip it
				//fmt.Printf("skipping %v\n", *table)
				continue
			} else {
				fmt.Printf("\"%v\"\n", *table)
			}
		}
	}
}

// Resource holds the type of AWS resource we are going to look for
var Resource string

func init() {
	rootCmd.AddCommand(orphanedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dryrunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	orphanedCmd.Flags().StringVarP(&Resource, "resource", "r", "AWS::DynamoDB::Table", "Which type of resource to enumerate")

}

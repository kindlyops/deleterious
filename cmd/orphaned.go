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

//go:generate mockery --name ".*API"

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// orphanedCmd represents the orphaned command.
var orphanedCmd = &cobra.Command{
	Use:   "orphaned",
	Short: "Find orphaned resources without deleting",
	Long:  `Finds all of the specified resources that are not referenced from an active CloudFormation stack.`,
	Run:   orphaned,
}

func orphaned(cmd *cobra.Command, args []string) {
	handler := func(rootedResources map[string]bool) {
		log.Fatal().Msgf("'%v' is not yet supported, supported types are %s", Resource, supportedTypes)
	}

	switch {
	case Resource == "AWS::DynamoDB::Table":
		handler = func(rootedResources map[string]bool) {
			svc := dynamodb.New(getSession())
			processDynamoDB(rootedResources, svc)
		}
	case Resource == "AWS::KMS::Key":
		handler = func(rootedResources map[string]bool) {
			svc := kms.New(getSession())
			processKMS(rootedResources, svc)
		}
	case Resource == "AWS::Kinesis::Stream":
		handler = func(rootedResources map[string]bool) {
			svc := kinesis.New(getSession())
			processKinesis(rootedResources, svc)
		}
	case Resource == "AWS::Logs::LogGroup":
		handler = func(rootedResources map[string]bool) {
			svc := cloudwatchlogs.New(getSession())
			l := lambda.New(getSession())
			processLogs(rootedResources, svc, l)
		}
	case Resource == "AWS::S3::Bucket":
		handler = func(rootedResources map[string]bool) {
			svc := s3.New(getSession())
			processS3(rootedResources, svc)
		}
	}

	cfn := cloudformation.New(getSession())
	rootedResources := getRootedResources(cfn, Resource)
	handler(rootedResources)
}

type CloudformationAPI interface {
	ListStacks(params *cloudformation.ListStacksInput) (*cloudformation.ListStacksOutput, error)
	ListStackResources(params *cloudformation.ListStackResourcesInput) (*cloudformation.ListStackResourcesOutput, error)
}

func getStackStates() []*string {
	stackStates := []*string{
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

	return stackStates
}

func getRootedResources(svc CloudformationAPI, kind string) map[string]bool {
	rootedResources := make(map[string]bool)

	var token *string

	for ok := true; ok; ok = (token != nil) {
		stacks, err := svc.ListStacks(&cloudformation.ListStacksInput{
			NextToken:         token,
			StackStatusFilter: getStackStates(),
		})
		if err != nil {
			log.Fatal().Msgf("Error listing stacks: %v", err)
		} else {
			token = stacks.NextToken
		}

		for _, stack := range stacks.StackSummaries {
			log.Debug().Msgf("Processing %v\n", *stack.StackName)

			var resourceToken *string

			for ok := true; ok; ok = (resourceToken != nil) {
				resources, err := svc.ListStackResources(&cloudformation.ListStackResourcesInput{
					StackName: stack.StackName,
					NextToken: resourceToken,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Error processing %v", *stack.StackName)

					break
				}

				resourceToken = resources.NextToken

				for _, resource := range resources.StackResourceSummaries {
					if *resource.ResourceType == kind {
						log.Debug().Msgf("Found rooted resource %v : %v\n", *resource.PhysicalResourceId, *resource.ResourceType)

						rootedResources[*resource.PhysicalResourceId] = true
					}
				}
			}
		}
	}

	return rootedResources
}

func getSession() *session.Session {
	// by default, AWS_PROFILE set in the env will control where the SDK
	// looks for credentials. If the --profile flag is specified, use that
	// to override the environment.
	if awsProfile != "" {
		os.Setenv("AWS_PROFILE", awsProfile)
	}

	return session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}))
}

type LogsAPI interface {
	DescribeLogGroups(params *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreams(params *cloudwatchlogs.DescribeLogStreamsInput) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
}

type LambdaAPI interface {
	GetFunction(params *lambda.GetFunctionInput) (*lambda.GetFunctionOutput, error)
}

func processLogs(rootedResources map[string]bool, svc LogsAPI, lambdaService LambdaAPI) {
	fmt.Fprintf(os.Stdout, "GroupName, RetentionDays, StoredBytesRaw, StoredBytesHuman, LastLogEntry, DaysAgo\n")

	var nextGroup *string

	const pageSize = 50

	for ok := true; ok; ok = (nextGroup != nil) {
		groups, err := svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: nextGroup,
			Limit:     aws.Int64(pageSize),
		})

		if err != nil {
			log.Error().Err(err).Msg("Error listing Log Groups")
		} else {
			if groups.NextToken != nil {
				log.Debug().Msg("Has more groups, setting nextGroup")
				nextGroup = groups.NextToken
			} else {
				nextGroup = nil
			}
		}

		for _, group := range groups.LogGroups {
			log.Debug().Msgf("Processing %v\n", *group)

			_, ownedByCfn := rootedResources[*group.LogGroupName]
			// check whether there is a lambda function with a matching name
			// we do this because log groups are freqently created by the
			// lambda, not by the cloudformation stack that created the lambda
			ownedByLambda := isLogGroupOwnedByLambda(group, lambdaService)

			if ownedByCfn || ownedByLambda {
				// this stream is owned by a cloudformation stack, skip it
				if ownedByCfn {
					log.Debug().Msgf("Group %v is owned by a cloudformation stack, skipping\n", *group.LogGroupName)
				}
			} else {
				// find when the last event was written
				lastEvent, daysAgo, retentionDays := getLogGroupLastWrites(svc, group)
				sizeHuman := humanize.Bytes(uint64(*group.StoredBytes))
				sizeRaw := uint64(*group.StoredBytes)
				fmt.Fprintf(os.Stdout, "\"%v\", %d, %v, %v, %v, %d\n",
					*group.LogGroupName, retentionDays, sizeRaw, sizeHuman, lastEvent, daysAgo)
			}
		}
	}
}

func getLogGroupLastWrites(svc LogsAPI, group *cloudwatchlogs.LogGroup) (
	lastEvent string, daysAgo int64, retentionDays int64) {
	stream, err := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(*group.LogGroupName),
		OrderBy:      aws.String("LastEventTime"),
		Descending:   aws.Bool(true),
		Limit:        aws.Int64(1),
	})

	lastEvent = "Never"

	const hoursInDay = 24

	if err != nil || stream == nil {
		log.Error().Err(err).Msg("Error listing Log Streams")

		lastEvent = "Unknown"
		daysAgo = -1
		retentionDays = -1
	} else if stream.LogStreams != nil &&
		len(stream.LogStreams) > 0 &&
		stream.LogStreams[0].LastEventTimestamp != nil {
		t := time.UnixMilli(*stream.LogStreams[0].LastEventTimestamp)
		lastEvent = t.Format(time.RFC3339)
		daysAgo = int64(time.Since(t).Hours() / hoursInDay)
	}

	if group.RetentionInDays != nil {
		retentionDays = *group.RetentionInDays
	}

	return lastEvent, daysAgo, retentionDays
}

func isLogGroupOwnedByLambda(group *cloudwatchlogs.LogGroup, lambdaService LambdaAPI) bool {
	ownedByLambda := false

	if strings.HasPrefix(*group.LogGroupName, "/aws/lambda") {
		lambdaName := strings.Split(*group.LogGroupName, "/aws/lambda/")[1]
		log.Debug().Msgf("Checking lambda %v", lambdaName)

		_, err := lambdaService.GetFunction(&lambda.GetFunctionInput{
			FunctionName: aws.String(lambdaName),
		})

		if err == nil {
			ownedByLambda = true

			log.Debug().Msgf("Group %v is owned by a lambda function, skipping", *group.LogGroupName)
		}
	}

	return ownedByLambda
}

type S3API interface {
	ListBuckets(params *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

func processS3(rootedResources map[string]bool, svc S3API) {
	buckets, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Error().Err(err).Msg("Error listing S3 Buckets")
	}

	var orphaned []s3.Bucket

	for _, bucket := range buckets.Buckets {
		log.Debug().Msgf("Processing %v\n", *bucket.Name)

		if _, ok := rootedResources[*bucket.Name]; ok {
			// this bucket is owned by a cloudformation stack, skip it
			log.Debug().Msgf("Bucket %v is owned by a cloudformation stack, skipping", *bucket)
		} else {
			orphaned = append(orphaned, *bucket)
		}
	}

	output, _ := json.MarshalIndent(orphaned, "", "  ")

	fmt.Fprintf(os.Stdout, "%s\n", string(output))
}

type KinesisAPI interface {
	ListStreams(params *kinesis.ListStreamsInput) (*kinesis.ListStreamsOutput, error)
}

func processKinesis(rootedResources map[string]bool, svc KinesisAPI) {
	var startStream *string

	const pageSize = 50

	for ok := true; ok; ok = (startStream != nil) {
		streams, err := svc.ListStreams(&kinesis.ListStreamsInput{
			ExclusiveStartStreamName: startStream,
			Limit:                    aws.Int64(pageSize),
		})

		if err != nil {
			log.Error().Err(err).Msg("Error listing Kinesis Streams")
		} else {
			if *streams.HasMoreStreams {
				log.Debug().Msg("Has more streams, setting startStream")
				startStream = streams.StreamNames[len(streams.StreamNames)-1]
			} else {
				startStream = nil
			}
		}

		for _, stream := range streams.StreamNames {
			log.Debug().Msgf("Processing %v", *stream)

			if _, ok := rootedResources[*stream]; ok {
				// this stream is owned by a cloudformation stack, skip it
				log.Debug().Msgf("Stream %v is owned by a cloudformation stack, skipping", *stream)
			} else {
				fmt.Fprintf(os.Stdout, "\"%v\"\n", *stream)
			}
		}
	}
}

type KmsAPI interface {
	ListKeys(params *kms.ListKeysInput) (*kms.ListKeysOutput, error)
	DescribeKey(params *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error)
}

func processKMS(rootedResources map[string]bool, svc KmsAPI) {
	var marker *string

	for ok := true; ok; ok = (marker != nil) {
		keys, err := svc.ListKeys(&kms.ListKeysInput{
			Marker: marker,
		})
		if err != nil {
			log.Error().Err(err).Msg("Error listing KMS keys")
		} else {
			marker = keys.NextMarker
		}

		for _, key := range keys.Keys {
			if _, ok := rootedResources[*key.KeyId]; ok {
				// this key is owned by a cloudformation stack, skip it
			} else {
				metadata, err := svc.DescribeKey(&kms.DescribeKeyInput{
					KeyId: key.KeyId,
				})
				if err != nil {
					// There are some strange key IDs that are returned by
					// ListKeys but that were never successfully created.
					// These keys don't exist, and if we can't DescribeKey,
					// treat it as not existing. If you are running with
					// reduced permissions you could also hit the same permission
					// denied error. but this is still the best indicator that
					// we have that a key is not there.
					continue
				}
				mgr := metadata.KeyMetadata.KeyManager
				if *mgr == "AWS" {
					// don't mess with AWS managed keys
					continue
				}
				state := metadata.KeyMetadata.KeyState
				if *state == "Enabled" {
					// candidate for cleanup
					fmt.Fprintf(os.Stdout, "\"%v\"\n", *key.KeyId)
				}
			}
		}
	}
}

type DynamoDBAPI interface {
	ListTables(params *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error)
}

func processDynamoDB(rootedResources map[string]bool, svc DynamoDBAPI) {
	var lastTable *string

	for ok := true; ok; ok = (lastTable != nil) {
		tables, err := svc.ListTables(&dynamodb.ListTablesInput{
			ExclusiveStartTableName: lastTable,
		})
		if err != nil {
			log.Error().Err(err).Msg("Error listing DynamoDB tables")
		} else {
			lastTable = tables.LastEvaluatedTableName
		}

		for _, table := range tables.TableNames {
			if _, ok := rootedResources[*table]; ok {
				// this table table belongs to a cloudformation stack, skip it
				log.Debug().Msgf("skipping %v\n", *table)

				continue
			} else {
				fmt.Fprintf(os.Stdout, "\"%v\"\n", *table)
			}
		}
	}
}

// Resource holds the type of AWS resource we are going to look for.
var Resource string

var supportedTypes = `
	AWS::DynamoDB::Table
	AWS::KMS::Key
	AWS::Kinesis::Stream
	AWS::Logs::LogGroup
	AWS::S3::Bucket`

func init() {
	rootCmd.AddCommand(orphanedCmd)

	orphanedCmd.Flags().StringVarP(&Resource, "resource", "r", "",
		fmt.Sprintf("Which type of resource to enumerate\nSupported types are%s", supportedTypes))

	_ = orphanedCmd.MarkFlagRequired("resource")
}

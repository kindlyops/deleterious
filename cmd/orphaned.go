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

	svc := cloudformation.New(
		session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
		})))

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

	stacks, err := svc.ListStacks(&cloudformation.ListStacksInput{
		StackStatusFilter: states,
	})
	if err != nil {
		fmt.Printf("Error listing stacks: %v", err)
	}
	//fmt.Printf("Got %v results\n", len(stacks.StackSummaries))
	rootedResources := make(map[string]bool)

	for _, stack := range stacks.StackSummaries {
		//fmt.Printf("Processing %v\n", *stack.StackName)
		resources, err := svc.ListStackResources(&cloudformation.ListStackResourcesInput{
			StackName: stack.StackName,
		})
		if err != nil {
			fmt.Printf("Error processing %v\n", *stack.StackName)
			continue
		}
		for _, resource := range resources.StackResourceSummaries {
			if *resource.ResourceType == Resource {
				//fmt.Printf("Found resource %v : %v\n", *resource.LogicalResourceId, *resource.ResourceType)
				rootedResources[*resource.PhysicalResourceId] = true
			}
		}
	}

	switch {
	case Resource == "AWS::DynamoDB::Table":
		processDynamoDB(rootedResources)
	default:
		fmt.Printf("%v is not yet supported", Resource)
	}
}

func processDynamoDB(rootedResources map[string]bool) {
	svc := dynamodb.New(
		session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
		})))
	tables, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		fmt.Printf("Error listing DynamoDB tables: %v\n", err)
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

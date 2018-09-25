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

	"github.com/spf13/cobra"
)

// dryrunCmd represents the dryrun command
var dryrunCmd = &cobra.Command{
	Use:   "dryrun",
	Short: "Find resources without deleting",
	Long:  `Finds all of the specified resources that are not referenced from an active CloudFormation stack.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dryrun called, resource is: %s\n", Resource)
	},
}

var Resource string

func init() {
	rootCmd.AddCommand(dryrunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dryrunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dryrunCmd.Flags().StringVarP(&Resource, "resource", "r", "AWS::DynamoDB::Table", "Which type of resource to enumerate")

}

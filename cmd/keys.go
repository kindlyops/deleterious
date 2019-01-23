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
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

// ExtendedAccessKeyMetadata extends the key metadata struct
type ExtendedAccessKeyMetadata struct {
	AccessKeyID string
	CreateDate  *time.Time
	LastUsed    *time.Time
	Status      string
	UserName    string
	Arn         string
	Age         int
}

// lastUsedCmd represents the lastUsed command
var lastUsedCmd = &cobra.Command{
	Use:   "apikeys",
	Short: "Find keys based on window of time in days.",
	Long:  `Finds all keys used given a window of time in days.`,
	Run:   lastUsed,
}

func lastUsed(cmd *cobra.Command, args []string) {

	svc := iam.New(getIamSession())
	listOfKeys := []ExtendedAccessKeyMetadata{}
	// TODO add pagination
	var maxItems int64 = 999
	users, err := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: &maxItems,
	})
	if err != nil {
		fmt.Printf("Error listing users: %v", err)
	}
	for _, user := range users.Users {
		keys, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
			MaxItems: &maxItems,
			UserName: user.UserName,
		})
		if err != nil {
			fmt.Printf("Error listing keys: %v", err)
		}
		for _, key := range keys.AccessKeyMetadata {
			used, err := svc.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
				AccessKeyId: key.AccessKeyId,
			})
			if err != nil {
				fmt.Printf("Error listing a key: %v", err)
			}

			lastUsedDate := used.AccessKeyLastUsed.LastUsedDate
			now := time.Now()
			// dateMarkerString is our current date minus days given to search for.
			dateMarkerString := now.AddDate(0, 0, -Days).Format(time.RFC3339)
			dateMarkerTime, err := time.Parse(time.RFC3339, dateMarkerString)
			if err != nil {
				fmt.Println("Error parsing the time.")
			}
			newKey := ExtendedAccessKeyMetadata{
				AccessKeyID: *key.AccessKeyId,
				CreateDate:  key.CreateDate,
				LastUsed:    lastUsedDate,
				Status:      *key.Status,
				UserName:    *key.UserName,
				Age:         int(time.Since(*key.CreateDate).Hours() / 24),
				Arn:         *user.Arn,
			}
			// dateTest is a bool which is true if the last used is after the marker time.
			dateTest := lastUsedDate.After(dateMarkerTime)
			// if Used = false we want to know what keys have not been used within x days.
			if Used == dateTest {
				listOfKeys = append(listOfKeys, newKey)
			}

		}
	}
	output, err := json.MarshalIndent(listOfKeys, "", "  ")
	fmt.Fprintf(os.Stdout, "%s", output)
}

func getIamSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}))
}

// Days holds the amount of days to search for.
var Days int

// Used is a bool to invert the search if needed.
var Used bool

func init() {
	rootCmd.AddCommand(lastUsedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dryrunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	lastUsedCmd.Flags().IntVar(&Days, "days", 30, "How many days to search for.")
	lastUsedCmd.Flags().BoolVar(&Used, "used", false, "Display only used keys in the last X days. (Defaults to false.)")

}

// Copyright © 2024 Kindly Ops, LLC <support@kindlyops.com>
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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// EcrAPI is the interface for ECR service calls used by this command.
type EcrAPI interface {
	DescribeRepositories(params *ecr.DescribeRepositoriesInput) (*ecr.DescribeRepositoriesOutput, error)
	ListImages(params *ecr.ListImagesInput) (*ecr.ListImagesOutput, error)
}

type repoImageSummary struct {
	RepositoryName string
	UntaggedCount  int
	TotalCount     int
}

var ecrCmd = &cobra.Command{
	Use:   "ecr",
	Short: "Report untagged images in ECR repositories",
	Long:  `Enumerates all ECR repositories and reports untagged vs total image counts as a markdown summary table.`,
	Run:   ecrReport,
}

func ecrReport(cmd *cobra.Command, args []string) {
	svc := ecr.New(getSession())
	processECR(svc)
}

func processECR(svc EcrAPI) {
	repos, err := listAllRepositories(svc)
	if err != nil {
		log.Fatal().Err(err).Msg("Error listing ECR repositories")
	}

	summaries := make([]repoImageSummary, 0, len(repos))

	for _, repo := range repos {
		repoName := aws.StringValue(repo.RepositoryName)
		log.Debug().Msgf("Processing repository %v", repoName)

		untagged, err := countImages(svc, repoName, "UNTAGGED")
		if err != nil {
			log.Error().Err(err).Msgf("Error counting untagged images for %v", repoName)

			continue
		}

		total, err := countImages(svc, repoName, "")
		if err != nil {
			log.Error().Err(err).Msgf("Error counting total images for %v", repoName)

			continue
		}

		summaries = append(summaries, repoImageSummary{
			RepositoryName: repoName,
			UntaggedCount:  untagged,
			TotalCount:     total,
		})
	}

	printMarkdownTable(summaries)
}

func listAllRepositories(svc EcrAPI) ([]*ecr.Repository, error) {
	var repos []*ecr.Repository

	var token *string

	for ok := true; ok; ok = (token != nil) {
		output, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
			NextToken: token,
		})
		if err != nil {
			return nil, err
		}

		repos = append(repos, output.Repositories...)
		token = output.NextToken
	}

	return repos, nil
}

func countImages(svc EcrAPI, repoName, tagStatus string) (int, error) {
	count := 0

	var token *string

	for ok := true; ok; ok = (token != nil) {
		input := &ecr.ListImagesInput{
			RepositoryName: aws.String(repoName),
			NextToken:      token,
		}

		if tagStatus != "" {
			input.Filter = &ecr.ListImagesFilter{
				TagStatus: aws.String(tagStatus),
			}
		}

		output, err := svc.ListImages(input)
		if err != nil {
			return 0, err
		}

		count += len(output.ImageIds)
		token = output.NextToken
	}

	return count, nil
}

func printMarkdownTable(summaries []repoImageSummary) {
	fmt.Fprintf(os.Stdout, "| Repository | Untagged Images | Total Images |\n")
	fmt.Fprintf(os.Stdout, "|---|---|---|\n")

	for _, s := range summaries {
		fmt.Fprintf(os.Stdout, "| %s | %d | %d |\n", s.RepositoryName, s.UntaggedCount, s.TotalCount)
	}
}

func init() {
	rootCmd.AddCommand(ecrCmd)
}

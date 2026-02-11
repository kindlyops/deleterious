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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/kindlyops/deleterious/cmd/mocks"
)

func captureStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	r.Close()

	return string(out)
}

func Test_listAllRepositories(t *testing.T) {
	t.Parallel()

	t.Run("single page", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{
					{RepositoryName: aws.String("my-app")},
					{RepositoryName: aws.String("my-api")},
				},
			}, nil).Once()

		repos, err := listAllRepositories(mockSvc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repos) != 2 {
			t.Errorf("expected 2 repos, got %d", len(repos))
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("paginated", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		nextToken := "page2"

		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{
					{RepositoryName: aws.String("repo-1")},
				},
				NextToken: aws.String(nextToken),
			}, nil).Once()

		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{
			NextToken: aws.String(nextToken),
		}).Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []*ecr.Repository{
				{RepositoryName: aws.String("repo-2")},
			},
		}, nil).Once()

		repos, err := listAllRepositories(mockSvc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repos) != 2 {
			t.Errorf("expected 2 repos, got %d", len(repos))
		}

		if aws.StringValue(repos[0].RepositoryName) != "repo-1" {
			t.Errorf("expected repo-1, got %s", aws.StringValue(repos[0].RepositoryName))
		}

		if aws.StringValue(repos[1].RepositoryName) != "repo-2" {
			t.Errorf("expected repo-2, got %s", aws.StringValue(repos[1].RepositoryName))
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("no repositories", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{},
			}, nil).Once()

		repos, err := listAllRepositories(mockSvc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repos) != 0 {
			t.Errorf("expected 0 repos, got %d", len(repos))
		}

		mockSvc.AssertExpectations(t)
	})
}

func Test_countImages(t *testing.T) {
	t.Parallel()

	t.Run("untagged images", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app"),
			Filter: &ecr.ListImagesFilter{
				TagStatus: aws.String("UNTAGGED"),
			},
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
				{ImageDigest: aws.String("sha256:bbb")},
			},
		}, nil).Once()

		count, err := countImages(mockSvc, "my-app", "UNTAGGED")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if count != 2 {
			t.Errorf("expected 2, got %d", count)
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("all images no filter", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app"),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
				{ImageDigest: aws.String("sha256:bbb")},
				{ImageDigest: aws.String("sha256:ccc")},
			},
		}, nil).Once()

		count, err := countImages(mockSvc, "my-app", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if count != 3 {
			t.Errorf("expected 3, got %d", count)
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("paginated images", func(t *testing.T) {
		t.Parallel()

		mockSvc := &mocks.EcrAPI{}
		nextToken := "page2"

		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app"),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
			},
			NextToken: aws.String(nextToken),
		}, nil).Once()

		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app"),
			NextToken:      aws.String(nextToken),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:bbb")},
				{ImageDigest: aws.String("sha256:ccc")},
			},
		}, nil).Once()

		count, err := countImages(mockSvc, "my-app", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if count != 3 {
			t.Errorf("expected 3, got %d", count)
		}

		mockSvc.AssertExpectations(t)
	})
}

func Test_processECR(t *testing.T) {
	t.Run("single repo with untagged images", func(t *testing.T) {
		mockSvc := &mocks.EcrAPI{}

		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{
					{RepositoryName: aws.String("my-app-backend")},
				},
			}, nil).Once()

		// untagged count
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app-backend"),
			Filter: &ecr.ListImagesFilter{
				TagStatus: aws.String("UNTAGGED"),
			},
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
				{ImageDigest: aws.String("sha256:bbb")},
			},
		}, nil).Once()

		// total count
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("my-app-backend"),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
				{ImageDigest: aws.String("sha256:bbb")},
				{ImageDigest: aws.String("sha256:ccc")},
				{ImageDigest: aws.String("sha256:ddd")},
				{ImageDigest: aws.String("sha256:eee")},
			},
		}, nil).Once()

		output := captureStdout(func() {
			processECR(mockSvc)
		})

		if !strings.Contains(output, "my-app-backend") {
			t.Errorf("expected output to contain repo name, got: %s", output)
		}

		if !strings.Contains(output, "| my-app-backend | 2 | 5 |") {
			t.Errorf("expected output to contain '| my-app-backend | 2 | 5 |', got: %s", output)
		}

		if !strings.Contains(output, "| Repository | Untagged Images | Total Images |") {
			t.Errorf("expected output to contain header row, got: %s", output)
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("no repositories", func(t *testing.T) {
		mockSvc := &mocks.EcrAPI{}

		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{},
			}, nil).Once()

		output := captureStdout(func() {
			processECR(mockSvc)
		})

		if !strings.Contains(output, "| Repository | Untagged Images | Total Images |") {
			t.Errorf("expected header row even with no repos, got: %s", output)
		}

		// Should only have header and separator, no data rows
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 2 {
			t.Errorf("expected 2 lines (header + separator), got %d: %s", len(lines), output)
		}

		mockSvc.AssertExpectations(t)
	})

	t.Run("multiple repos", func(t *testing.T) {
		mockSvc := &mocks.EcrAPI{}

		mockSvc.On("DescribeRepositories", &ecr.DescribeRepositoriesInput{}).
			Return(&ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{
					{RepositoryName: aws.String("backend")},
					{RepositoryName: aws.String("frontend")},
				},
			}, nil).Once()

		// backend untagged
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("backend"),
			Filter: &ecr.ListImagesFilter{
				TagStatus: aws.String("UNTAGGED"),
			},
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
			},
		}, nil).Once()

		// backend total
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("backend"),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:aaa")},
				{ImageDigest: aws.String("sha256:bbb")},
				{ImageDigest: aws.String("sha256:ccc")},
			},
		}, nil).Once()

		// frontend untagged
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("frontend"),
			Filter: &ecr.ListImagesFilter{
				TagStatus: aws.String("UNTAGGED"),
			},
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{},
		}, nil).Once()

		// frontend total
		mockSvc.On("ListImages", &ecr.ListImagesInput{
			RepositoryName: aws.String("frontend"),
		}).Return(&ecr.ListImagesOutput{
			ImageIds: []*ecr.ImageIdentifier{
				{ImageDigest: aws.String("sha256:ddd")},
			},
		}, nil).Once()

		output := captureStdout(func() {
			processECR(mockSvc)
		})

		if !strings.Contains(output, "| backend | 1 | 3 |") {
			t.Errorf("expected backend row, got: %s", output)
		}

		if !strings.Contains(output, "| frontend | 0 | 1 |") {
			t.Errorf("expected frontend row, got: %s", output)
		}

		mockSvc.AssertExpectations(t)
	})
}

func Test_printMarkdownTable(t *testing.T) {
	summaries := []repoImageSummary{
		{RepositoryName: "my-app", UntaggedCount: 15, TotalCount: 42},
		{RepositoryName: "my-api", UntaggedCount: 0, TotalCount: 8},
	}

	output := captureStdout(func() {
		printMarkdownTable(summaries)
	})

	expected := []string{
		"| Repository | Untagged Images | Total Images |",
		"|---|---|---|",
		"| my-app | 15 | 42 |",
		"| my-api | 0 | 8 |",
	}

	for _, line := range expected {
		if !strings.Contains(output, line) {
			t.Errorf("expected output to contain %q, got: %s", line, output)
		}
	}
}

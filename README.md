# deleterious
Helps you clean up AWS resources. This is handy when
retention policies on CloudFormation stacks leave lots
of orphaned AWS resources around costing money.

## installation for homebrew (MacOS/Linux)

    brew install kindlyops/tap/deleterious

once installed, you can upgrade to a newer version using this command:

    brew upgrade kindlyops/tap/deleterious

## installation for scoop (Windows Powershell)

To enable the bucket for your scoop installation

    scoop bucket add kindlyops https://github.com/kindlyops/kindlyops-scoop
    
To install deleterious

    scoop install deleterious

once installed, you can upgrade to a newer version using this command:

    scoop status
    scoop update deleterious

## installation from source

    go get github.com/kindlyops/deleterious
    deleterious help

## Example of deleting DynamoDB tables

Once deleterious gives you a list of things to delete, and
you have manually confirmed they are ok to delete, you
can make a little loop to delete the objects. Here is an example with dynamoDB tables

```bash
#!/bin/bash

# tables that need to be deleted
declare -a tables=("foo-MonkeyTable-1FDTVGZJOT25Y"
"foo-BananaTable-1HFLQZL7CVQ7L"
)

for i in "${tables[@]}"; do
	echo "deleting table: $i"
	aws dynamodb delete-table --table-name "$i"
done
```

## Example of deleting S3 buckets

To list orphaned buckets, you can run

    deleterious orphaned --resource "AWS::S3::Bucket"

To distill the output into only the bucket names, you can use jq

    deleterious orphaned --resource "AWS::S3::Bucket" | jq -r '.[].Name'

Once deleterious gives you a list of things to delete, and
you have manually confirmed they are ok to delete, you
can make a little loop to delete the objects. Here is an example with S3 buckets

```bash
#!/bin/bash

# buckets that need to be deleted
declare -a buckets=("foo-bananabucket-148lv5q85e3dc"
	"foo-bananabucket-14bh2oapj6a3e"
)

for i in "${buckets[@]}"; do
	echo "deleting bucket: $i"
	aws s3api delete-bucket --bucket "$i"
done
```

## Testing release process

To run goreleaser locally to test changes to the release process configuration:

    goreleaser release --snapshot --skip-publish --rm-dist

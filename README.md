# deleterious
cleans up AWS resources

## Example of deleting s3 buckets

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

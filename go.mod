module github.com/kindlyops/deleterious

go 1.12

require (
	github.com/aws/aws-sdk-go v1.30.8
	// If changing rules_go version, remember to change version in WORKSPACE also
	github.com/bazelbuild/rules_go v0.22.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.6.3
)

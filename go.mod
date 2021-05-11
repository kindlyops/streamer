module github.com/kindlyops/streamer

go 1.16

require (
	github.com/a8m/kinesis-producer v0.2.0
	github.com/aws/aws-sdk-go v1.38.38
	// If changing rules_go version, remember to change version in WORKSPACE also
	github.com/bazelbuild/rules_go v0.27.0
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0 // indirect
)

[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/kindlyops/streamer)

# utilities for working with kinesis streams

## streamer - kinesis pry bar

An utility for working with kinesis streams.

* to get a set of records that have failed ingest in a Kinesis Firehose,
  use the aws CLI to copy the records to your filesystem. For example:
  `aws s3 ls --summarize --human-readable --recursive 's3://my-bucket/firehose_errors/format-conversion-failed/'`
  and then
  `aws s3 cp --recursive s3://my-bucket/firehose_errors/format_conversion_failed/ error_data/`
* Extract will extract the original source record from a tree of error logs,
  filtering by error type. The records are decoded and batched into a new file.
* OPTIONAL: perform any needed processing/transform of the records using custom
  code to fix them up so ingest will succeed.
* Once the records are ready to attempt ingestion again, `load` will batch load
  the records into a specified kinesis stream, using the `kinesis-producer`
  library to control the ingest rate.

## installation for homebrew (MacOS/Linux)

    brew install kindlyops/tap/streamer

once installed, you can upgrade to a newer version using this command:

    brew upgrade kindlyops/tap/streamer

## installation for scoop (Windows Powershell)

To enable the bucket for your scoop installation

    scoop bucket add kindlyops https://github.com/kindlyops/kindlyops-scoop

To install deleterious

    scoop install streamer

once installed, you can upgrade to a newer version using this command:

    scoop status
    scoop update streamer

## installation from source

    go get github.com/kindlyops/streamer
    streamer help

## Developer instructions

Want to help add features or fix bugs? Awesome! streamer is built using bazel.

    `brew install bazelisk`
    grab the source code from github
    `bazel run streamer` to compile and run the locally compiled version

### Testing release process

To run goreleaser locally to test changes to the release process configuration:

    goreleaser release --snapshot --skip-publish --rm-dist

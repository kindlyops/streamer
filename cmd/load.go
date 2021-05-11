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
	"log"
	"os"
	"time"

	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load <directory> <stream>",
	Short: "Walk a tree of jsonl files and batch load into kinesis stream.",
	Long:  `Load a set of json records from a tree of files and load them into kinesis in batches with backpressure.`,
	Run:   load,
	Args:  cobra.ExactArgs(2), // TODO
}

func load(cmd *cobra.Command, args []string) {
	data := args[0]
	target := args[1]
	_, _ = os.Stat(data)

	client := kinesis.New(session.New(aws.NewConfig()))
	pr := producer.New(&producer.Config{
		StreamName:   target,
		BacklogCount: 2000,
		Client:       client,
	})

	pr.Start()

	// Handle failures
	go func() {
		for r := range pr.NotifyFailures() {
			// r contains `Data`, `PartitionKey` and `Error()`
			log.Print(r)
		}
	}()

	go func() {
		for i := 0; i < 5000; i++ {
			err := pr.Put([]byte("foo"), "bar")
			if err != nil {
				log.Fatalf("error producing %v", err)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	pr.Stop()

	log.Fatal("TODO")
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

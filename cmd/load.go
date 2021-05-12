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
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	producer "github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load <directory> <stream>",
	Short: "Walk a tree of jsonl files and batch load into kinesis stream.",
	Long:  `Load a set of json records from a tree of files and load them into kinesis in batches with backpressure.`,
	Run:   load,
	Args:  cobra.ExactArgs(2), //nolint:gomnd // this is an appropriate magic number
}

func load(cmd *cobra.Command, args []string) {
	data := args[0]
	target := args[1]

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Error establishing AWS session: %v", err)
	}

	client := kinesis.New(sess)
	pr := producer.New(&producer.Config{
		StreamName:   target,
		BacklogCount: 2000, //nolint:gomnd // channel backlog before blocking Put
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

	visit := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("failed to open file %v: %w", p, err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if Debug {
					log.Printf("Got: %v\n", line)
				}

				err = pr.Put([]byte(line), "test_partition_key")
				if err != nil {
					log.Fatalf("Error producing to kinesis: %v\n", err)
				}
			}
		}

		return nil
	}

	err = filepath.Walk(data, visit)
	if err != nil {
		log.Printf("Error from filepath.Walk: %v\n", err)
	}

	pr.Stop()
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

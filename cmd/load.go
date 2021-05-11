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

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load <directory> <stream>",
	Short: "Walk a tree of jsonl files and batch load into kinesis stream.",
	Long:  `Load a set of json records from a tree of files and load them into kinesis in batches with backpressure.`,
	Run:   load,
	Args:  cobra.ExactArgs(1), // TODO
}

func load(cmd *cobra.Command, args []string) {
	target := args[0]
	_, _ = os.Stat(target)

	log.Fatal("TODO")
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

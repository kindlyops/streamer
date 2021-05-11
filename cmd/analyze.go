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

var analyzeCmd = &cobra.Command{
	Use:   "analyze <directory>",
	Short: "Analyze a tree of kinesis error logs.",
	Long:  `Load a set of json records from a tree of files and summarize the kinesis error codes.`,
	Run:   analyze,
	Args:  cobra.ExactArgs(1),
}

func analyze(cmd *cobra.Command, args []string) {
	target := args[0]
	_, _ = os.Stat(target)

	log.Fatal("TODO")
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

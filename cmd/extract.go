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

var extractCmd = &cobra.Command{
	Use:   "extract <directory> <output>",
	Short: "Walk a tree of kinesis error logs and extract the original records.",
	Long:  `Extract original records from a set of kinesis error logs and prepare for reprocessing.`,
	Run:   load,
	Args:  cobra.ExactArgs(2), // TODO
}

func extract(cmd *cobra.Command, args []string) {
	source := args[0]
	output := args[1]
	_, _ = os.Stat(source)
	_, _ = os.Stat(output)

	log.Fatal("TODO")
}

func init() {
	rootCmd.AddCommand(extractCmd)
}

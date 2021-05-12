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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract <errorlogs_dir> <output_dir>",
	Short: "Walk a tree of kinesis error logs and extract the original records from `rawData`.",
	Long:  `Extract original records from a set of kinesis error logs and prepare for reprocessing.`,
	Run:   extract,
	Args:  cobra.ExactArgs(2), //nolint:gomnd // this is an appropriate magic number
}

func extract(cmd *cobra.Command, args []string) {

	source := args[0]
	output := args[1]

	visit := func(p string, info os.FileInfo, err2 error) error {
		if err2 != nil {
			return err2
		}

		target := filepath.Join(output, strings.TrimPrefix(p, source))

		if !info.IsDir() {
			f, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("failed to open file %v: %w", p, err)
			}
			defer f.Close()

			out, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create %v: %w", target, err)
			}
			defer out.Close()

			w := bufio.NewWriter(out)

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if Debug {
					log.Printf("Got: %v\n", line)
				}

				raw, err := getOriginalRecord(line)
				if err != nil {
					log.Printf("Failed to extract original record from: %v\n", line)
				} else {
					_, err = w.Write(raw)
					if err != nil {
						return fmt.Errorf("failed to write to buffer %v: %w", target, err)
					}
					err = w.WriteByte('\n')
					if err != nil {
						return fmt.Errorf("failed to write to buffer %v: %w", target, err)
					}
				}
			}
			w.Flush()
		} else {
			err := os.MkdirAll(target, 0755)
			if err != nil {
				return fmt.Errorf("failed to create %v: %w", target, err)
			}
		}

		return nil
	}

	err := filepath.Walk(source, visit)
	if err != nil {
		log.Printf("Error from filepath.Walk: %v\n", err)
	}
}

type kinesisError struct {
	RawData       string `json:"rawdata"`
	LastErrorCode string `json:"lastErrorCode"`
}

func getOriginalRecord(line string) ([]byte, error) {
	data := kinesisError{}

	err := json.Unmarshal([]byte(line), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %v: %w", line, err)
	}

	decoded, err := b64.StdEncoding.DecodeString(data.RawData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %v: %w", data.RawData, err)
	}

	return decoded, nil
}

func init() {
	rootCmd.AddCommand(extractCmd)
}

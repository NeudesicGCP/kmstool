// Copyright 2018 Neudesic LLC
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
	"context"

	"github.com/NeudesicGCP/kmstool/io"
	"github.com/NeudesicGCP/kmstool/kms"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file using KMS",
		RunE:  encrypt,
	}
)

func init() {
	rootCmd.AddCommand(encryptCmd)
}

func encrypt(cmd *cobra.Command, args []string) error {
	Logger.Debug("Starting encrypt")
	ctx := context.Background()
	source, err := io.Read(ctx, plaintext)
	if err != nil {
		return errors.Wrap(err, "Error reading plaintext")
	}
	target, err := kms.Encrypt(ctx, key, source)
	if err != nil {
		return errors.Wrap(err, "Error encrypting plaintext")
	}
	if err = io.Write(ctx, ciphertext, target); err != nil {
		return errors.Wrap(err, "Error writing ciphertext")
	}
	Logger.Debug("Exiting encrypt")
	return err
}

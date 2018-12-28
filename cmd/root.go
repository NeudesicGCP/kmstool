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

// This package defines the commands to encrypt and decrypt
package cmd

import (
	"log"

	"github.com/blendle/zapdriver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// This is the structured logger that will write to stdout. It is
	// configured to be format compatible with Stackdriver, but it is not
	// responsible for delivering the output to Stackdriver.
	Logger     *zap.SugaredLogger
	key        string
	ciphertext string
	plaintext  string
	rootCmd    = &cobra.Command{
		Use:   "kmstool",
		Short: "kmstool is used to encrypt and decrypt files using Cloud KMS keys",
		Long:  `kmstool is a utility to encrypt and/or decrypt a file using a symmetric KMS key, that can read/write to a mounted file path, or a GCS bucket.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate that required parameters are present
			ciphertext = viper.GetString("ciphertext")
			if ciphertext == "" {
				return errors.New("ciphertext path must be provided")
			}
			key = viper.GetString("key")
			if key == "" {
				return errors.New("key name must be provided")
			}
			plaintext = viper.GetString("plaintext")
			if plaintext == "" {
				return errors.New("plaintext path must be provided")
			}
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(bootstrap)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to a TOML or YAML configuration file")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().StringP("key", "k", "", "Fully qualified KMS key name, in format project/PROJECTID/locations/LOCATION/keyRings/KEYRING/NAME")
	_ = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	_ = rootCmd.MarkFlagRequired("key")
	rootCmd.PersistentFlags().StringP("ciphertext", "s", "", "Path to ciphertext file; can be local path or gs://bucket/object")
	_ = viper.BindPFlag("ciphertext", rootCmd.PersistentFlags().Lookup("ciphertext"))
	_ = rootCmd.MarkFlagRequired("ciphertext")
	rootCmd.PersistentFlags().StringP("plaintext", "P", "", "Path to plaintext file; can be local path or gs://bucket/object")
	_ = viper.BindPFlag("plaintext", rootCmd.PersistentFlags().Lookup("plaintext"))
	_ = rootCmd.MarkFlagRequired("plaintext")
}

func bootstrap() {
	// Load configuration from file, flags, and environment
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Error locating home directory: %v", err)
	}
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigName(".kmstool")
	viper.SetEnvPrefix("kmstool")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	// Logger is dependent on verbosity, which could have come from flag, env,
	// or configuration file.
	Logger = buildLogger()

	if err == nil {
		return
	}

	switch t := err.(type) {
	case viper.ConfigFileNotFoundError:
		Logger.Debugw("Error reading configuration file",
			"error", t,
		)

	default:
		Logger.Errorw("Error reading configuration file",
			"error", t,
		)
	}
}

// Returns a sugared zap logger that meets the Stackdriver structured format
func buildLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	if viper.GetBool("verbose") {
		logger, err = zapdriver.NewDevelopment()
	} else {
		logger, err = zapdriver.NewProduction()
	}
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
		return nil
	}
	_ = zap.RedirectStdLog(logger)
	return logger.Sugar()
}

// Executes the command with supporting parameters
func Execute(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

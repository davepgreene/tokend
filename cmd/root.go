// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	http "github.com/davepgreene/tokend/http"
	"github.com/davepgreene/tokend/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// TokendCmd represents the base command when called without any subcommands
var TokendCmd = &cobra.Command{
	Use:   "tokend",
	Short: "Tokend interfaces with Hashicorp's Vault to provide a secure method to deliver secrets to servers in the cloud.",
	Long: `Tokend gives security and accountability around the delivery
	of secrets to servers running in the cloud.

	It provides a seamless interface between Vault and Rapid7's
	Propsd allowing developers to specify the secrets they need
	for their service without putting unencrypted secrets out in the wild.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initializeConfig()
		initializeLog()
		if err != nil {
			return err
		}

		return boot()
	},
}

func boot() error {
	router := http.Handler()
	log.Error(router)
	return router
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the TokendCmd.
func Execute() {
	if err := TokendCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	TokendCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	TokendCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose level logging")
	validConfigFilenames := []string{"json"}
	TokendCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

func initializeLog() {
	log.RegisterExitHandler(func() {
		log.Info("Shutting down")
	})

	// Set logging options based on config
	if lvl, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		log.SetLevel(lvl)
	} else {
		log.Info("Unable to parse log level in settings. Defaulting to INFO")
	}

	// If using verbose mode, log at debug level
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if cfgFile != "" {
		log.WithFields(log.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("Loaded config file")
	}

}

func initializeConfig(subCmdVs ...*cobra.Command) error {
	utils.Defaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	viper.AutomaticEnv() // read in environment variables that match

	return nil
}

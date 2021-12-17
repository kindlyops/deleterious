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
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var awsProfile string

// Debug controls whether or not debug messages should be printed.
var Debug bool

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Version: "dev",
	Use:     "deleterious",
	Short:   "Deletes stuff from AWS",
	Long: `Deleterious helps delete resources from AWS accounts
in order to reduce costs from unused items. Typically these
are resources from cloudformation stacks that have a deletion
policy set to retain, and iterting in a test account leaves
behind many orphaned databases or other costly items.
Brought to you by

_  ___           _ _        ___
| |/ (_)_ __   __| | |_   _ / _ \ _ __  ___
| ' /| | '_ \ / _| | | | | | | | | '_ \/ __|
| . \| | | | | (_| | | |_| | |_| | |_) __ \
|_|\_\_|_| |_|\__,_|_|\__, |\___/| .__/|___/
                      |___/      |_|
use at your own risk.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	rootCmd.SetVersionTemplate(v)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339

	if isatty.IsTerminal(os.Stdout.Fd()) {
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.With().Caller().Logger().Output(output)
	} else {
		log.Logger = log.With().Caller().Logger()
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.deleterious.yaml)")
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "AWS profile to use")

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Print debug messages while working")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error().Stack().Err(err).Msg("Couldn't locate home dir")
			os.Exit(1)
		}

		// Search config in home directory with name ".deleterious" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".deleterious")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msg(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}

	if Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

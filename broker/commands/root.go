//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package commands

import (
	"errors"
	"fmt"
	"os"

	api "github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "sentel",
	Short: "sentel is tool to dianose sentel server",
	Long:  `sentel can diagnose sentel servier status`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var (
	cfgFile   string
	brokerApi *api.BrokerApi
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file ")
	RootCmd.PersistentFlags().StringP("author", "a", "cloudstone", "cloudstone")
	RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "cloudstone")
	viper.SetDefault("license", "apache")

	// Version
	RootCmd.AddCommand(versionCmd)
	// Status
	RootCmd.AddCommand(statusCmd)
	// Service
	RootCmd.AddCommand(servicesCmd)

	// Clients
	clientsCmd.AddCommand(clientsShowCmd)
	clientsCmd.AddCommand(clientsKickoffCmd)
	RootCmd.AddCommand(clientsCmd)
	// Sessions
	sessionsCmd.AddCommand(sessionsShowCmd)
	RootCmd.AddCommand(sessionsCmd)
	// Subscriptions
	subscriptionsCmd.AddCommand(subscriptionsShowCmd)
	RootCmd.AddCommand(subscriptionsCmd)
	// Topics
	topicsCmd.AddCommand(topicsShowCmd)
	RootCmd.AddCommand(topicsCmd)

}

func Run(c config.Config) error {
	api, err := api.NewBrokerApi(c)
	if err != nil || api == nil {
		return errors.New("Sentel service is not started, please start sentel at first")
	}
	brokerApi = api
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return err
	}
	return nil
}

func initConfig() {
}

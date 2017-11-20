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

package cmd

import (
	"errors"

	api "github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "sentel-ctl",
	Short: "sentel-ctl is tool to dianose sentel server",
	Long:  `sentel-ctl can diagnose sentel servier status`,
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

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(brokerCmd)
	RootCmd.AddCommand(clientsCmd)
	RootCmd.AddCommand(clusterCmd)
	RootCmd.AddCommand(pluginsCmd)
	RootCmd.AddCommand(routesCmd)
	RootCmd.AddCommand(servicesCmd)
	RootCmd.AddCommand(sessionsCmd)
	RootCmd.AddCommand(subscriptionsCmd)
	RootCmd.AddCommand(topicsCmd)

}

func Execute() error {
	c, _ := core.NewConfigWithFile(cfgFile)
	api, err := api.NewBrokerApi(c)
	if err != nil {
		return errors.New("Sentel service is not started, please start sentel at first")
	}
	brokerApi = api
	return RootCmd.Execute()
}

func initConfig() {
}

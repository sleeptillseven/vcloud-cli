package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/floriankammermann/vcloud-cli/vcdapi"
	"github.com/spf13/viper"
)

var networkname string

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "execute queries",
	Long: `execute queries`,
}

var allocatedipsCmd = &cobra.Command{
	Use:   "allocatedips",
	Short: "allocatedips for an org network",
	Long: `get all allocated ips of an org network`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(networkname) > 0 {
			user := viper.GetString("user")
			password := viper.GetString("password")
			org := viper.GetString("org")
			vcdapi.GetAuthToken(user, password, org)
			vcdapi.GetAllocatedIpsForNetworkName(networkname)
		} else {
			fmt.Println("you have to provide the networkname")
		}
	},
}

func init() {
	queryCmd.AddCommand(allocatedipsCmd)
	allocatedipsCmd.Flags().StringVarP(&networkname, "network", "n", "", "network name to search allocated ips on")
	RootCmd.AddCommand(queryCmd)
}

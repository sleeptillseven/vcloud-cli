package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vcloud-cli",
	Short: "a command line interface for the vcloud director api",
	Long: `a command line interface for the vcloud director api`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	RootCmd.PersistentFlags().String("user", "", "Port to run Application server on")
	RootCmd.PersistentFlags().String( "password", "", "password of vcloud director api")
	RootCmd.PersistentFlags().String("org", "", "org of vcloud director api")
	viper.BindPFlag("user", RootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("org", RootCmd.PersistentFlags().Lookup("org"))

	viper.SetEnvPrefix("vcd") // will be uppercased automatically
	viper.AutomaticEnv()
}
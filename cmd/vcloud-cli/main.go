package main

import (
	"fmt"
	"os"
	"github.com/floriankammermann/vcloud-cli/vcdapi"
	"github.com/urfave/cli"
)

const vcdUserEnvVarName string = "VCD_USER"
const vcdPasswordEnvVarName string = "VCD_PASSWORD"
const vcdOrgEnvVarName string = "VCD_ORG"

func main() {

	var networkname string
	var user string
	var password string
	var org string

	app := cli.NewApp()
	app.Name = "vcloud-cli"
	app.Usage = "use one of the options"
	app.Version = "v0.1.0"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: vcdUserEnvVarName,
			Value: "",
			Usage: "the api user for the vcloud directory api",
			EnvVar: vcdUserEnvVarName,
			Destination: &user,
		},
		cli.StringFlag{
			Name: vcdPasswordEnvVarName,
			Value: "",
			Usage: "the api user password for the vcloud directory api",
			EnvVar: vcdPasswordEnvVarName,
			Destination: &password,
		},
		cli.StringFlag{
			Name: vcdOrgEnvVarName,
			Value: "",
			Usage: "the organisation for the vcloud directory api",
			EnvVar: vcdOrgEnvVarName,
			Destination: &org,
		},
		cli.StringFlag{
			Name: "networkname",
			Value: "",
			Usage: "provide the organisation networkname",
			Destination: &networkname,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "query",
			Aliases:     []string{"t"},
			Usage:       "query for resources over the vcloud director api",
			Subcommands: []cli.Command{
				{
					Name:  "allocatedips",
					Usage: "search for allocated ips",
					Action: func(c *cli.Context) error {
						if len(networkname) > 0 {
							vcdapi.GetAuthToken(user, password, org)
							vcdapi.GetAllocatedIpsForNetworkName(networkname)
						} else {
							fmt.Println("you have to provide the networkname")
						}
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)

	/*
	createCommand := flag.NewFlagSet("query", flag.ExitOnError)
	queryType := createCommand.String("type", "", "query type.")
	networkname := createCommand.String("networkname", "", "networkname")

	if len(os.Args) == 0 {
		fmt.Printf("you have to provide a valid command")
	}
	switch os.Args[1] {
	case "query":
		createCommand.Parse(os.Args[2:])
		query(*queryType, *networkname)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	*/
}
/*
func query(queryType string, networkname string) {
	fmt.Printf("query type: [%s]\n", queryType)
	vcdapi.GetAuthToken()
	if "vm" == queryType {
		vcdapi.GetAllVm()
	}
	if "allocatedip" == queryType {
		vcdapi.GetAllocatedIpForNetworkName(networkname)
	}
}
*/
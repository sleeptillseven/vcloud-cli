package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/floriankammermann/vcloud-cli/types"
	"github.com/juju/errors"
	"io/ioutil"
	"net/http"
	"os"
)

var vcdClient *VcdClientType

const vcdUserEnvVarName string = "VCD_USER"
const vcdPasswordEnvVarName string = "VCD_PASSWORD"
const vcdOrgEnvVarName string = "VCD_ORG"

type VcdClientType struct {
	VAToken string // vCloud Air authorization token
}

func main() {
	createCommand := flag.NewFlagSet("query", flag.ExitOnError)
	queryType := createCommand.String("type", "", "query type.")
	networkname := createCommand.String("networkname", "", "networkname")

	switch os.Args[1] {
	case "query":
		createCommand.Parse(os.Args[2:])
		query(*queryType, *networkname)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func getAuthToken() {
	user := os.Getenv(vcdUserEnvVarName)
	if len(user) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n", vcdUserEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdUserEnvVarName, user)

	password := os.Getenv(vcdPasswordEnvVarName)
	if len(password) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n", vcdPasswordEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdPasswordEnvVarName, "***************")

	org := os.Getenv(vcdOrgEnvVarName)
	if len(org) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n", vcdOrgEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdOrgEnvVarName, org)

	req, err := http.NewRequest("POST", "https://datacenter.swisscomcloud.com/api/sessions", nil)
	req.Header.Set("Accept", "application/*+xml;version=5.5")
	req.SetBasicAuth(user+"@"+org, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	auth := resp.Header.Get("x-vcloud-authorization")
	vcdClient = &VcdClientType{
		VAToken: auth,
	}
	fmt.Printf("authorization: [%s]\n", auth)
}

func query(queryType string, networkname string) {
	fmt.Printf("query type: [%s]\n", queryType)
	getAuthToken()
	if "vm" == queryType {
		getAllVm()
	}
	if "allocatedip" == queryType {
		getAllocatedIpForNetworkName(networkname)
	}
}

func getAllVm() {
	req, err := http.NewRequest("GET", "https://datacenter.swisscomcloud.com/api/query?type=vm&fields=name&pageSize=512", nil)
	req.Header.Set("x-vcloud-authorization", vcdClient.VAToken)
	req.Header.Set("Accept", "application/*+xml;version=5.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	queryRes := new(types.QueryResultRecordsType)
	decodeBody(resp, queryRes)

	for _, vm := range queryRes.VMRecord {
		fmt.Printf("VM Name [%s]\n", vm.Name)
	}

}

func getAllocatedIpForNetworkName(networkname string) error {

	if len(networkname) == 0 {
		return errors.New("networkname is empty")
	}
	req, err := http.NewRequest("GET", "https://datacenter.swisscomcloud.com/api/query?type=orgNetwork&fields=name&filter=name=="+networkname, nil)
	req.Header.Set("x-vcloud-authorization", vcdClient.VAToken)
	req.Header.Set("Accept", "application/*+xml;version=5.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	queryRes := new(types.QueryResultRecordsType)
	decodeBody(resp, queryRes)

	for _, net := range queryRes.OrgNetworkRecord {
		fmt.Printf("Org Network Name [%s]\n", net.Name)
	}

	if len(queryRes.OrgNetworkRecord) > 1 {
		return errors.New("found more than one org network for name: " + networkname)
	}
	if len(queryRes.OrgNetworkRecord) == 0 {
		return errors.New("found no org network for name: " + networkname)
	}
	getAllocatedIpsForNetworkHref(queryRes.OrgNetworkRecord[0].HREF)
	return nil
}

func getAllocatedIpsForNetworkHref(networkref string) error {
	if len(networkref) == 0 {
		return errors.New("networkref is empty")
	}

	fmt.Printf("the network href: [%s]\n", networkref)

	req, err := http.NewRequest("GET", networkref+"/allocatedAddresses", nil)
	req.Header.Set("x-vcloud-authorization", vcdClient.VAToken)
	req.Header.Set("Accept", "application/*+xml;version=5.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	queryRes := new(types.AllocatedIpAddressesType)
	decodeBody(resp, queryRes)

	for _, ipAddress := range queryRes.IpAddress {
		fmt.Printf("ip [%s] ", ipAddress.IpAddress)
		for _, link := range ipAddress.Link {
			if "down" == link.Rel && "application/vnd.vmware.vcloud.vApp+xml" == link.Type {
				fmt.Printf("vApp: [%s] ", link.Name)
			}
			if "down" == link.Rel && "application/vnd.vmware.vcloud.vm+xml" == link.Type {
				fmt.Printf("vm: [%s] ", link.Name)
			}
		}
		fmt.Print("\n")
	}

	return nil
}

// decodeBody is used to XML decode a response body
func decodeBody(resp *http.Response, out interface{}) error {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//fmt.Println("response Body:", string(body))

	// Unmarshal the XML.
	if err = xml.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

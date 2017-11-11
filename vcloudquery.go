package main

import (
    "fmt"
    "os"
	"net/http"
	"io/ioutil"
	"flag"
	"encoding/xml"
	"github.com/ukcloud/govcloudair/types/v56"
)

var vcdClient *VcdClientType
const vcdUserEnvVarName string = "VCD_USER"
const vcdPasswordEnvVarName string = "VCD_PASSWORD"
const vcdOrgEnvVarName string = "VCD_ORG"

type VcdClientType struct {
	VAToken       string  // vCloud Air authorization token
}

func main() {
	queryType := flag.String("query", "", "query type.")
	flag.Parse()
	fmt.Printf("query type: [%s]\n", *queryType)
	getAuthToken()
	if "vm" == *queryType {
		getAllVm()
	}

}

func getAuthToken() {
	user := os.Getenv(vcdUserEnvVarName)
	if len(user) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n",vcdUserEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdUserEnvVarName, user)

	password := os.Getenv(vcdPasswordEnvVarName)
	if len(password) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n",vcdPasswordEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdPasswordEnvVarName, "***************")

	org := os.Getenv(vcdOrgEnvVarName)
	if len(org) == 0 {
		fmt.Printf("the environment variable [%s] is not set\n",vcdOrgEnvVarName)
		os.Exit(-1)
	}
	fmt.Printf("%s: [%s]\n", vcdOrgEnvVarName, org)

	req, err := http.NewRequest("POST", "https://datacenter.swisscomcloud.com/api/sessions", nil)
	req.Header.Set("Accept", "application/*+xml;version=1.5")
	req.SetBasicAuth(user + "@" + org, password)

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
	fmt.Printf("authorization: [%s]", vcdClient.VAToken)
}

func getAllVm() {
	req, err := http.NewRequest("GET", "https://datacenter.swisscomcloud.com/api/query?type=vm&fields=name&pageSize=512", nil)
	req.Header.Set("x-vcloud-authorization", vcdClient.VAToken)
	req.Header.Set("Accept", "application/*+xml;version=1.5")

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

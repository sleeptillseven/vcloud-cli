package main

import (
    "fmt"
    "os"
	"net/http"
	"io/ioutil"
	"flag"
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

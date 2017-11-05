package main

import (
    "fmt"
    "os"
	"net/http"
	"io/ioutil"
)

var vcdClient *VcdClientType

type VcdClientType struct {
	VAToken       string  // vCloud Air authorization token
}

func main() {

	user := os.Getenv("VCD_USER")
    fmt.Println("VCD_USER:", user)

    password := os.Getenv("VCD_PASSWORD")
    fmt.Println("VCD_PASSWORD:", "")

    org := os.Getenv("VCD_ORG")
	fmt.Println("VCD_ORG:", org)

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

	getAllVm()
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

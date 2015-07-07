package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const IP_URL = "https://ip-ranges.amazonaws.com/ip-ranges.json"

type IpInfo struct {
	Ip_prefix, Region, Service string
}

type Response struct {
	SyncToken, CreateDate string
	Prefixes              []IpInfo
}

func main() {
	delPtr := flag.String("delimiter", "\n", "delimiter")
	flag.Parse()

	r := getJson()
	ips := make(chan string)
	done := make(chan bool)
	// Totally unnecessary but fun :)
	go getIps(r.Prefixes, ips)
	go printIps(ips, done, *delPtr)
	<-done
}

func getJson() Response {
	resp, err := http.Get(IP_URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func printIps(ips chan string, done chan bool, delimiter string) {
	first := true
	for {
		ip, more := <-ips
		if more {
			if first {
				fmt.Print(ip)
				first = false
			} else {
				fmt.Printf("%s%s", delimiter, ip)
			}
		} else {
			break
		}
	}
	done <- true
	close(done)
}

func getIps(prefixes []IpInfo, ips chan string) {
	for _, prefix := range prefixes {
		if prefix.Service == "AMAZON" {
			ips <- prefix.Ip_prefix
		}
	}
	close(ips)
}

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
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
	r := getJson()
	ips := make(chan string)
	done := make(chan bool)
	// Totally unnecessary but fun :)
	go getIps(r.Prefixes, ips)
	go printIps(ips, done)
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

func printIps(ips chan string, done chan bool){
	first := true
	for {
		ip, more := <-ips
		if more {
			if( first ){
				fmt.Print(ip)
				first = false
			} else {
				fmt.Print(",", ip)
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

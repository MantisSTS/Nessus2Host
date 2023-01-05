package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type NessusClientDataV2 struct {
	Report []struct {
		Name       string `xml:"name,attr"`
		ReportHost []struct {
			Text           string `xml:",chardata"`
			Name           string `xml:"name,attr"`
			HostProperties struct {
				Text string `xml:",chardata"`
				Tag  []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"tag"`
			} `xml:"HostProperties"`
			ReportItem []struct {
				Port string `xml:"port,attr"`
			} `xml:"ReportItem"`
		} `xml:"ReportHost"`
	} `xml:"Report"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {

	// Use flags to specify the nessus xml file
	file := flag.String("f", "", "Nessus XML file")
	flag.Parse()

	// Read the file
	if *file == "" {
		fmt.Println("Please specify a file")
		os.Exit(1)
	}

	// Read the nessus xml file
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Unmarshal the xml into a Report struct
	var report NessusClientDataV2
	err = xml.Unmarshal(data, &report)
	if err != nil {
		fmt.Println("Error unmarshalling xml:", err)
		os.Exit(1)
	}

	var hosts []string

	// Get the host-ip and port
	for _, report := range report.Report {
		for _, host := range report.ReportHost {
			for _, item := range host.ReportItem {
				if item.Port == "0" {
					continue
				}
				host := fmt.Sprintf("%s:%s", host.Name, item.Port)
				// Check if hosts contains the host, if not then append it
				if !contains(hosts, host) {
					hosts = append(hosts, host)
				}
			}
		}
	}

	for _, host := range hosts {
		fmt.Println(host)
	}
}

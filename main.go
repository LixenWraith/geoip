package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// ipQuery is the request format required by ip-api.com's batch endpoint.
type ipQuery struct {
	Query string `json:"query"`
}

// ipResponse represents the JSON response from ip-api.com for each IP.
type ipResponse struct {
	Status  string `json:"status"`            // "success" or "fail"
	Country string `json:"country,omitempty"` // present if status=="success"
	Message string `json:"message,omitempty"` // error message if status=="fail"
}

const (
	batchSize = 100 // ip-api.com supports up to 100 queries per batch request
	apiURL    = "http://ip-api.com/batch?fields=status,country,message"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <ip_list_file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	ips, err := readAndValidateIPs(filename)
	if err != nil {
		log.Fatalf("Error reading IP file: %v", err)
	}

	if len(ips) == 0 {
		log.Fatalf("No valid IP addresses found in file.")
	}

	countryCounts := make(map[string]int)

	// Process the IPs in batches
	for i := 0; i < len(ips); i += batchSize {
		end := i + batchSize
		if end > len(ips) {
			end = len(ips)
		}
		batch := ips[i:end]

		results, err := lookupBatch(batch)
		if err != nil {
			log.Printf("Error looking up batch starting at index %d: %v", i, err)
			continue
		}

		for _, res := range results {
			if res.Status != "success" {
				log.Printf("Lookup failed for an IP: %s", res.Message)
				continue
			}
			countryCounts[res.Country]++
		}

		// Optional: if you are making many requests, consider sleeping between batches.
		time.Sleep(1500 * time.Millisecond)
	}

	// Print out the country counts
	for country, count := range countryCounts {
		fmt.Printf("%s: %d\n", country, count)
	}
}

// readAndValidateIPs reads a file, trims each line, and returns a slice of valid IPv4 address strings.
func readAndValidateIPs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Validate the IP using net.ParseIP.
		if netIP := net.ParseIP(line); netIP == nil {
			log.Printf("Line %d: '%s' is not a valid IP address; skipping.", lineNum, line)
			continue
		}
		ips = append(ips, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ips, nil
}

// lookupBatch queries the ip-api.com batch endpoint for a slice of IP addresses.
func lookupBatch(ips []string) ([]ipResponse, error) {
	// Prepare the JSON payload for the batch request.
	queries := make([]ipQuery, len(ips))
	for i, ip := range ips {
		queries[i] = ipQuery{Query: ip}
	}
	payload, err := json.Marshal(queries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch payload: %v", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("HTTP POST error: %v", err)
	}
	defer resp.Body.Close()

	var results []ipResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return results, nil
}

package main

import (
	"context"
	"os"

	"github.com/SealJonny/dyndns-go/server"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
)

func listDNSRecords(client *cloudflare.Client, zoneID string) ([]dns.RecordResponse, error) {
	page, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return nil, err
	}

	return page.Result, nil
}

func main() {
	token, exists := os.LookupEnv("CF_TOKEN")
	if !exists {
		panic("CF_TOKEN is not set")
	}

	zoneID, exists := os.LookupEnv("CF_ZONE_ID")
	if !exists {
		panic("CF_ZONE_ID is not set")
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	server := server.New(port, token, zoneID)
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

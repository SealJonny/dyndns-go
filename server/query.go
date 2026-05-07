package server

import (
	"fmt"
	"net"
	"net/url"
)

type QueryArgs struct {
	accountID string
	token     string
	domain    string
	ipv4      string
}

func NewQueryArgs(query url.Values) (*QueryArgs, error) {
	if !query.Has("domain") || !query.Has("ipv4") || !query.Has("token") {
		return nil, fmt.Errorf("must specify domain, ipv4 and token query arguments")
	}

	accountID := query.Get("accountID")
	if len(accountID) == 0 {
		return nil, fmt.Errorf("empty accountID")
	}

	token := query.Get("token")
	if len(token) == 0 {
		return nil, fmt.Errorf("empty token")
	}

	domain := query.Get("domain")
	if len(domain) == 0 {
		return nil, fmt.Errorf("invalid domain")
	}

	ipv4 := query.Get("ipv4")
	ip := net.ParseIP(ipv4)
	if ip == nil || ip.To4() == nil {
		return nil, fmt.Errorf("invalid ipv4 address")
	}

	return &QueryArgs{accountID, token, domain, ipv4}, nil
}

package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/SealJonny/dyndns-go/cloudflare"
)

type Server struct {
	port       string
	cloudflare cloudflare.CloudflareClient
}

func New(port string, token string, zoneID string) *Server {
	return &Server{
		port,
		*cloudflare.New(token, zoneID),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/dyndns", s.handleDynDNS)
	return http.ListenAndServe(":"+s.port, mux)
}

func (s *Server) handleDynDNS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	if !query.Has("domain") || !query.Has("ipv4") {
		http.Error(w, "must specify domain and ipv4 query arguments", http.StatusBadRequest)
		return
	}

	domain := query.Get("domain")
	if _, err := net.LookupHost(domain); err != nil {
		http.Error(w, "invalid domain", http.StatusBadRequest)
		return
	}

	ipv4 := query.Get("ipv4")
	ip := net.ParseIP(ipv4)
	if ip == nil || ip.To4() == nil {
		http.Error(w, "invalid ipv4 address", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	record, err := s.cloudflare.GetARecordForDomain(ctx, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedRecord, err := s.cloudflare.UpdateARecord(ctx, record, ipv4)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, updatedRecord.JSON.RawJSON())
}

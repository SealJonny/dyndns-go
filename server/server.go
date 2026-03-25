package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SealJonny/dyndns-go/cf"
	"github.com/SealJonny/dyndns-go/notification"
)

type Server struct {
	port   string
	zoneID string
}

func New(port string, zoneID string) *Server {
	return &Server{
		port,
		zoneID,
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

	args, err := NewQueryArgs(query)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	client := cf.New(args.accountID, args.token, s.zoneID)
	if err := client.VerifyToken(ctx); err != nil {
		slog.Warn("invalid token or account id")
		http.Error(w, "invalid token or account id", http.StatusBadRequest)
		return
	}

	record, err := client.GetARecordByDomain(ctx, args.domain)
	if err != nil {
		slog.Error("failed to list dns records for domain", "err", err, "domain", args.domain)
		notification.SMTPError("Cloudflare List Error", "Could not list dns records for domain.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	records, err := client.GetARecordsByIPv4(ctx, record.Content)
	if err != nil {
		slog.Error("failed to list dns records for ipv4", "err", err, "ipv4", record.Content)
		notification.SMTPError("Cloudflare List Error", "Could not list dns records for ipv4")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, r := range records {
		_, err = client.UpdateARecord(ctx, &r, args.ipv4)
		if err != nil {
			slog.Error("failed to update A record", "domain", r.Name, "ipv4", args.ipv4, "err", err)
			notification.SMTPError("Cloudflare Update Error", fmt.Sprintf("Could not update %s", r.Name))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		slog.Info("successfully updated", "domain", r.Name)
	}

	notification.SMTPInfo("Updated DNS records", fmt.Sprintf("Updated DNS record for %s successfully.", args.domain))
}

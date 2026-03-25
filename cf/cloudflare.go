package cf

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/accounts"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

type CloudflareDNS struct {
	zoneID    string
	accountID string
	client    *cloudflare.Client
}

func New(accountID string, token string, zoneID string) *CloudflareDNS {
	return &CloudflareDNS{
		zoneID:    zoneID,
		accountID: accountID,
		client:    cloudflare.NewClient(option.WithAPIToken(token)),
	}
}

func (c *CloudflareDNS) VerifyToken(ctx context.Context) error {
	response, err := c.client.Accounts.Tokens.Verify(ctx, accounts.TokenVerifyParams{
		AccountID: cloudflare.F(c.accountID),
	})
	if err != nil {
		return err
	}

	if response.Status != accounts.TokenVerifyResponseStatusActive {
		return fmt.Errorf("token must be active: token is %s", response.Status)
	}

	return nil
}

// GetARecordByDomain returns the matching A record for the given domain.
func (c *CloudflareDNS) GetARecordByDomain(ctx context.Context, domain string) (*dns.RecordResponse, error) {
	page, err := c.client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(c.zoneID),
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
		Name:   cloudflare.F(dns.RecordListParamsName{Exact: cloudflare.F(domain)}),
	})
	if len(page.Result) != 1 {
		return nil, fmt.Errorf("there should be only one A records for %s: %d record(s)", domain, len(page.Result))
	}

	return &page.Result[0], err
}

// GetARecordsByIPv4 returns all A records in the zone that match the given IPv4 address.
// It paginates through all results returned by the Cloudflare API.
func (c *CloudflareDNS) GetARecordsByIPv4(ctx context.Context, ipv4 string) ([]dns.RecordResponse, error) {
	page, err := c.client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID:  cloudflare.F(c.zoneID),
		Type:    cloudflare.F(dns.RecordListParamsTypeA),
		Content: cloudflare.F(dns.RecordListParamsContent{Exact: cloudflare.F(ipv4)}),
	})
	if err != nil {
		return nil, err
	}

	var records []dns.RecordResponse

	for page != nil {
		records = append(records, page.Result...)

		page, err = page.GetNextPage()
		if err != nil {
			return nil, err
		}
	}

	return records, err
}

// UpdateARecord updates the given record with the provided IPv4 address and returns it.
func (c *CloudflareDNS) UpdateARecord(ctx context.Context, record *dns.RecordResponse, ipv4 string) (*dns.RecordResponse, error) {
	return c.client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body: dns.ARecordParam{
			Content: cloudflare.F(ipv4),
		},
	})
}

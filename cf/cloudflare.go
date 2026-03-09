package cf

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/accounts"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

type CloudflareClient struct {
	zoneID    string
	accountID string
	client    *cloudflare.Client
}

func New(accountID string, token string, zoneID string) *CloudflareClient {
	return &CloudflareClient{
		zoneID:    zoneID,
		accountID: accountID,
		client:    cloudflare.NewClient(option.WithAPIToken(token)),
	}
}

func (c *CloudflareClient) VerifyToken(ctx context.Context) error {
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

func (c *CloudflareClient) GetARecordForDomain(ctx context.Context, domain string) (*dns.RecordResponse, error) {
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

func (c *CloudflareClient) UpdateARecord(ctx context.Context, record *dns.RecordResponse, ipv4 string) (*dns.RecordResponse, error) {
	return c.client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body: dns.ARecordParam{
			Content: cloudflare.F(ipv4),
		},
	})
}

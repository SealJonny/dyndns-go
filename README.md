# dyndns-go

A lightweight Dynamic DNS server written in Go. It exposes a simple HTTP endpoint that, when called, updates a Cloudflare DNS A record with the provided IPv4 address — making it easy to integrate with routers or scripts that support custom DynDNS providers.

## How It Works

`dyndns-go` runs an HTTP server and listens for `GET /dyndns` requests. On each request it:

1. Validates the query parameters
2. Verifies the provided Cloudflare API token
3. Looks up the existing A record for the given domain in its Cloudflare zone
4. Updates that A record with the new IPv4 address

## Requirements

- A [Cloudflare](https://cloudflare.com) account with the target domain managed as a zone
- A Cloudflare **account** API token (not a user token) with **DNS Edit** permissions for the zone
- Go 1.22+ (only needed to build from source)

## Configuration

The server is configured via environment variables:

| Variable      | Required | Default | Description                                      |
|---------------|----------|---------|--------------------------------------------------|
| `CF_ZONE_ID`  | ✅ Yes   | —       | The Cloudflare Zone ID of the targeted domain    |
| `PORT`        | ❌ No    | `80`    | The port the HTTP server listens on              |

## API

### `GET /dyndns`

Updates the Cloudflare DNS A record for a domain.

**Query Parameters:**

| Parameter   | Required | Description                                      |
|-------------|----------|--------------------------------------------------|
| `accountID` | ✅ Yes   | Your Cloudflare account ID                       |
| `token`     | ✅ Yes   | A Cloudflare **account** API token (not a user token) with DNS Edit permissions |
| `domain`    | ✅ Yes   | The fully qualified domain name to update        |
| `ipv4`      | ✅ Yes   | The new IPv4 address to set on the A record      |

**Example:**

```
GET /dyndns?accountID=abc123&token=mytoken&domain=home.example.com&ipv4=1.2.3.4
```

**Responses:**

| Status | Meaning                                         |
|--------|-------------------------------------------------|
| `200`  | A record updated successfully                   |
| `400`  | Missing/invalid parameters or invalid token     |
| `405`  | Method not allowed (only GET is supported)      |
| `500`  | Cloudflare API error                            |

## Installation

### Download a Release

Pre-built binaries are available on the [Releases](../../releases) page.

### Build from Source

```bash
git clone https://github.com/SealJonny/dyndns-go.git
cd dyndns-go
go build -o bin/dyndns-go ./...
```

## Usage

```bash
export CF_ZONE_ID=your_zone_id
export PORT=8080  # optional

./bin/dyndns-go
```

The server will log startup information and listen for requests on the configured port.

## Router Integration

Most home routers with a custom DynDNS option (e.g. Fritz!Box, pfSense) can be pointed at this server. Configure the router to call:

```
http://<server-host>:<port>/dyndns?accountID=<id>&token=<token>&domain=<domain>&ipv4=<ipaddr>
```

Many routers will substitute the current WAN IP automatically into a placeholder like `<ipaddr>`.

# For what?
I have AdGuard Home on my PC. When I turn pc off, I need to switch DNS server on router (for me I need to turn off 'strict order' in OpenWRT settings).

I dont wanna go to settings every night and every morning.

# What is that?
So I built simple DNS server which has my AdGuard Home as first upstream and Cloudflare as second upstream server. And when my PC is down it uses Cloudflare. When my PC is up it uses AdGuard Home on it. Automatically. That's what I needed.

# How to use


> [!NOTE]
> You have to create config.json to run server (see [Config](#config))

## Standalone binary
- Download binary for your platform from releases tab
- Rename to dns-switcher (if desired)
- Run: `./dns-switcher`

## Docker

### CLI
`docker run -p 53:53 -v ./config.json:/app/config.json -d ghcr.io/misshanya/dns-switcher`

### Docker Compose
```yaml
services:
  dns-server:
    container_name: dns-switcher-server
    image: ghcr.io/misshanya/dns-switcher
    ports:
      - "53:53"
    volumes:
      - ./config.json:/app/config.json
    restart: unless-stopped
```

# Config
You should configure it via `config.json`

You can configure listen address and upstream servers. Example config:
```json
{
    "address": ":53",
    "upstreams": [
        "1.1.1.1:53",
        "1.0.0.1:53"
    ]
}
```

# For what?
I have AdGuard Home on my PC. When I turn pc off, I need to switch DNS server on router (for me I need to turn off 'strict order' in OpenWRT settings).

I dont wanna go to settings every night and every morning.

# What is that?
So I built simple DNS server which has my AdGuard Home as first upstream and Cloudflare as second upstream server. And when my PC is down it uses Cloudflare. When my PC is up it uses AdGuard Home on it. Automatically. That's what I needed.

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

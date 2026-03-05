# sandbox-egress-header

A sandbox template that runs a transparent proxy ([mitmproxy](https://mitmproxy.org/)) to inject a custom `X-Sandbox-ID` header into all outbound HTTP and HTTPS requests.

## Template ID

`sandbox-egress-header`

## What's Pre-installed

- **mitmproxy** — Transparent proxy running in the background (mitmdump)
- **iptables rules** — Redirect all outbound HTTP/HTTPS traffic through the proxy
- **System-wide CA certificate** — mitmproxy's CA is trusted so HTTPS interception works transparently
- **Header injection** — Every outbound request gets an `X-Sandbox-ID` header set to the sandbox ID

## How It Works

1. On sandbox start, a transparent [mitmproxy](https://mitmproxy.org/) proxy starts on port **18181**
2. iptables rules redirect all outbound HTTP (80) and HTTPS (443) traffic through the proxy
3. The proxy injects the sandbox ID as a header into every outbound request
4. A compiled `sockmark.so` library (loaded via `LD_PRELOAD`) marks the proxy's own sockets so iptables skips them, preventing redirect loops

## Configuration

| Environment variable | Default | Description |
|---|---|---|
| `HEADER_NAME` | `X-Sandbox-ID` | Name of the header injected into outbound requests |
| `PROXY_PORT` | `18181` | Port the transparent proxy listens on |

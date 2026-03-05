# sandbox-egress-header

A sandbox template that runs a transparent proxy ([mitmproxy](https://mitmproxy.org/)) to inject a custom `X-Sandbox-ID` header into all outbound HTTP and HTTPS requests.

## Template ID

`sandbox-egress-header`

## What's Pre-installed

- **mitmproxy** — Transparent proxy running in the background (mitmdump)
- **iptables rules** — Redirect all outbound HTTP/HTTPS traffic through the proxy
- **System-wide CA certificate** — mitmproxy's CA is trusted so HTTPS interception works transparently
- **Header injection** — Every outbound request gets an `X-Sandbox-ID` header set to the sandbox ID

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("sandbox-egress-header", timeout=120)
try:
    # All outbound HTTP/HTTPS requests automatically include X-Sandbox-ID
    result = sbx.commands.run("curl -s https://httpbin.org/headers")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('sandbox-egress-header', { timeoutMs: 120_000 });
try {
  // All outbound HTTP/HTTPS requests automatically include X-Sandbox-ID
  const result = await sbx.commands.run('curl -s https://httpbin.org/headers');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

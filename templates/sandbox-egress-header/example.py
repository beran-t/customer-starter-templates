import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"sandbox-egress-header:{tag}" if tag else "sandbox-egress-header"

sbx = Sandbox.create(template_ref, timeout=120)
try:
    # Verify mitmproxy is running
    ps = sbx.commands.run("ps aux | grep mitmdump | grep -v grep")
    assert "mitmdump" in ps.stdout, "mitmdump not found in process list"

    # Verify X-Sandbox-ID header is injected into HTTP requests
    http = sbx.commands.run("curl -s http://httpbin.org/headers", timeout=30)
    assert "X-Sandbox-Id" in http.stdout or "X-Sandbox-ID" in http.stdout, (
        f"X-Sandbox-ID header not found in HTTP response: {http.stdout}"
    )

    # Verify X-Sandbox-ID header is injected into HTTPS requests
    https = sbx.commands.run(
        "curl -s https://httpbin.org/headers",
        timeout=30,
        envs={"SSL_CERT_FILE": "/etc/ssl/certs/ca-certificates.crt"},
    )
    assert "X-Sandbox-Id" in https.stdout or "X-Sandbox-ID" in https.stdout, (
        f"X-Sandbox-ID header not found in HTTPS response: {https.stdout}"
    )

    print("All checks passed.")
finally:
    sbx.kill()

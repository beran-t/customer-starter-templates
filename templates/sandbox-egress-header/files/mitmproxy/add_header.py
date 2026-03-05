"""
mitmproxy addon that injects X-Sandbox-ID header into every outbound HTTP/HTTPS request.

Used with mitmdump in transparent mode. The sandbox ID is read from
/run/e2b/.E2B_SANDBOX_ID (written by E2B's envd at sandbox boot).

The file may be empty briefly after sandbox restore (envd writes it
asynchronously), so we retry on every request until a real value appears.
"""

import os
from mitmproxy import http

HEADER_NAME = os.environ.get("HEADER_NAME", "X-Sandbox-ID")
SANDBOX_ID_FILE = "/run/e2b/.E2B_SANDBOX_ID"
_sandbox_id: str | None = None


def request(flow: http.HTTPFlow) -> None:
    global _sandbox_id

    if _sandbox_id is not None:
        flow.request.headers[HEADER_NAME] = _sandbox_id
        return

    # Try reading the sandbox ID file (written by envd after sandbox boot).
    sid = ""
    try:
        with open(SANDBOX_ID_FILE) as f:
            sid = f.read().strip()
    except Exception:
        pass

    # Fall back to E2B_SANDBOX_ID env var.
    if not sid or sid == "unknown":
        sid = os.environ.get("E2B_SANDBOX_ID", "").strip()

    # Cache only once we have a real value; otherwise retry next request.
    if sid and sid != "unknown":
        _sandbox_id = sid
    else:
        sid = "unknown"

    flow.request.headers[HEADER_NAME] = sid

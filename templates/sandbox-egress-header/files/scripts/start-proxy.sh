#!/usr/bin/env bash
set -euo pipefail

PROXY_PORT="${PROXY_PORT:-18181}"
ADDON_PATH="/opt/mitmproxy/add_header.py"
SOCKMARK_LIB="/opt/mitmproxy/sockmark.so"
CONFDIR="/root/.mitmproxy"

echo "=== Egress Header Proxy Setup ==="
echo "Proxy port: ${PROXY_PORT}"

# --- 0. Resolve sandbox ID from MMDS (Firecracker metadata service) ---
# The E2B_SANDBOX_ID env var isn't available in the start command's process,
# so we fetch it directly from the instance metadata service.
MMDS_TOKEN=$(curl -sf -X PUT "http://169.254.169.254/latest/api/token" \
    -H "X-metadata-token-ttl-seconds: 21600" 2>/dev/null || echo "")
if [ -n "${MMDS_TOKEN}" ]; then
    SANDBOX_ID=$(curl -sf -H "X-metadata-token: ${MMDS_TOKEN}" \
        "http://169.254.169.254/instanceID" 2>/dev/null || echo "unknown")
else
    SANDBOX_ID="${E2B_SANDBOX_ID:-unknown}"
fi
echo "Sandbox ID: ${SANDBOX_ID}"

# --- 1. Generate mitmproxy CA certificate ---
echo "Generating mitmproxy CA certificate..."
mitmdump --mode transparent --listen-port "${PROXY_PORT}" \
    --set confdir="${CONFDIR}" \
    -s "${ADDON_PATH}" &
GEN_PID=$!

for i in $(seq 1 30); do
    [ -f "${CONFDIR}/mitmproxy-ca-cert.pem" ] && break
    sleep 0.5
done

kill "${GEN_PID}" 2>/dev/null || true
wait "${GEN_PID}" 2>/dev/null || true

# --- 2. Install CA certificate system-wide ---
echo "Installing CA certificate..."
cp "${CONFDIR}/mitmproxy-ca-cert.pem" /usr/local/share/ca-certificates/mitmproxy.crt
update-ca-certificates 2>/dev/null

# Set CA env vars for Python, Node.js, and OpenSSL
cat > /etc/profile.d/mitmproxy-ca.sh << 'ENVEOF'
export REQUESTS_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt
export SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
export NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt
ENVEOF

cat >> /etc/environment << 'ENVEOF'
REQUESTS_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt
SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt
ENVEOF

# --- 3. Set up iptables rules ---
# Skip packets marked with 1 (mitmproxy's own outbound traffic via LD_PRELOAD + SO_MARK)
echo "Configuring iptables..."
iptables -t nat -A OUTPUT -m mark --mark 1 -j RETURN
iptables -t nat -A OUTPUT -d 127.0.0.0/8 -j RETURN
iptables -t nat -A OUTPUT -p tcp --dport 80 -j REDIRECT --to-port "${PROXY_PORT}"
iptables -t nat -A OUTPUT -p tcp --dport 443 -j REDIRECT --to-port "${PROXY_PORT}"

# --- 4. Start mitmproxy with LD_PRELOAD for loop prevention ---
echo "Starting mitmproxy in transparent mode..."
LD_PRELOAD="${SOCKMARK_LIB}" \
    E2B_SANDBOX_ID="${SANDBOX_ID}" \
    HEADER_NAME="${HEADER_NAME:-X-Sandbox-ID}" \
    mitmdump --mode transparent --listen-port "${PROXY_PORT}" \
    -s "${ADDON_PATH}" --set confdir="${CONFDIR}" \
    --quiet &
MITM_PID=$!

# --- 5. Wait for proxy to be ready ---
sleep 2
if kill -0 "${MITM_PID}" 2>/dev/null; then
    echo "Proxy is ready (PID: ${MITM_PID})"
else
    echo "WARNING: mitmproxy may not have started correctly"
fi
touch /tmp/proxy-ready

echo "=== Egress Header Proxy Running ==="
wait "${MITM_PID}"

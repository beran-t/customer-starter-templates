import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `sandbox-egress-header:${tag}` : 'sandbox-egress-header';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 120_000 });
try {
  // Verify mitmproxy is running
  const ps = await sbx.commands.run('ps aux | grep mitmdump | grep -v grep');
  if (!ps.stdout.includes('mitmdump')) throw new Error('mitmdump not found in process list');

  // Verify X-Sandbox-ID header is injected into HTTP requests
  const http = await sbx.commands.run('curl -s http://httpbin.org/headers', { timeoutMs: 30_000 });
  if (!http.stdout.includes('X-Sandbox-Id') && !http.stdout.includes('X-Sandbox-ID')) {
    throw new Error(`X-Sandbox-ID header not found in HTTP response: ${http.stdout}`);
  }

  // Verify X-Sandbox-ID header is injected into HTTPS requests
  const https = await sbx.commands.run('curl -s https://httpbin.org/headers', {
    timeoutMs: 30_000,
    envs: { SSL_CERT_FILE: '/etc/ssl/certs/ca-certificates.crt' },
  });
  if (!https.stdout.includes('X-Sandbox-Id') && !https.stdout.includes('X-Sandbox-ID')) {
    throw new Error(`X-Sandbox-ID header not found in HTTPS response: ${https.stdout}`);
  }

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

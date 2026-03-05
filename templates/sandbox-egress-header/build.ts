import "dotenv/config"
import { Template, defaultBuildLogger, waitForFile } from 'e2b'

const TEMPLATE_NAME = 'sandbox-egress-header'
const tag = process.env.E2B_BUILD_TAG || undefined

export const template = Template()
  .fromBaseImage()
  .aptInstall([
    'python3',
    'python3-pip',
    'iptables',
    'ca-certificates',
    'curl',
    'gcc',
    'libc-dev',
  ])
  .runCmd('pip3 install mitmproxy --break-system-packages', { user: 'root' })
  .runCmd('mkdir -p /opt/mitmproxy', { user: 'root' })
  .copy('files/mitmproxy/add_header.py', '/opt/mitmproxy/add_header.py')
  .copy('files/mitmproxy/sockmark.c', '/opt/mitmproxy/sockmark.c')
  .copy('files/scripts/start-proxy.sh', '/opt/mitmproxy/start-proxy.sh')
  .runCmd('gcc -shared -fPIC -o /opt/mitmproxy/sockmark.so /opt/mitmproxy/sockmark.c -ldl', { user: 'root' })
  .runCmd('chmod +x /opt/mitmproxy/start-proxy.sh', { user: 'root' })
  .setStartCmd('sudo /opt/mitmproxy/start-proxy.sh', waitForFile('/tmp/proxy-ready'))

Template.build(template, tag ? `${TEMPLATE_NAME}:${tag}` : TEMPLATE_NAME, {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

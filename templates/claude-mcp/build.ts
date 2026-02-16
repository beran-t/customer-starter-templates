import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromUbuntuImage("25.04")

  // --- Base system deps ---
  .runCmd([
    "sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release build-essential",
  ])

  // --- Docker CE ---
  .runCmd([
    "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg",
    'echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null',
    "sudo apt-get update && sudo apt-get install -y docker-ce docker-ce-cli containerd.io jq",
    "sudo mkdir -p /root/.docker",
    'echo \'{"max-concurrent-downloads":10,"max-download-attempts":10}\' | sudo tee /etc/docker/daemon.json > /dev/null',
    'echo \'{"HttpHeaders":{"User-Agent":"E2B/0.0.1"}}\' | sudo tee /root/.docker/config.json > /dev/null',
    "sudo apt-get clean && sudo rm -rf /var/lib/apt/lists/*",
  ])

  // --- Python, pip, git ---
  .runCmd([
    "sudo apt-get update && sudo apt-get install -y python3 python3-pip python3-venv git && sudo apt-get clean && sudo rm -rf /var/lib/apt/lists/*",
  ])

  // --- Poetry ---
  .runCmd([
    "sudo pip install --break-system-packages pipx && pipx install poetry && sudo cp ~/.local/bin/poetry /usr/local/bin/",
  ])

  // --- Node.js 22 ---
  .runCmd([
    "sudo mkdir -p /etc/apt/keyrings && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && echo \"deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main\" | sudo tee /etc/apt/sources.list.d/nodesource.list > /dev/null",
    "sudo apt-get update && sudo apt-get install -y --no-install-recommends nodejs && sudo corepack enable || true",
  ])

  // --- uv / uvx ---
  .runCmd([
    "curl -LsSf https://astral.sh/uv/install.sh | sh && sudo cp ~/.local/bin/uv /usr/local/bin/ && sudo cp ~/.local/bin/uvx /usr/local/bin/",
  ])

  // --- Claude Code CLI ---
  .runCmd([
    "curl -fsSL https://claude.ai/install.sh | bash && sudo cp -L ~/.local/bin/claude /usr/local/bin/claude",
  ])

  // --- Go 1.25 (needed to compile gateway) ---
  .runCmd([
    "curl -fsSL https://go.dev/dl/go1.25.7.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -",
  ])

  // --- Copy gateway source and compile ---
  .copy("./mcp-gateway-src", "/opt/mcp-gateway-src")
  .runCmd([
    "cp -r /opt/mcp-gateway-src /tmp/mcp-build && cd /tmp/mcp-build && export PATH=$PATH:/usr/local/go/bin && go build -o /tmp/type-gen ./cmd/type-gen && /tmp/type-gen && CGO_ENABLED=0 go build -o /tmp/mcp-gateway ./cmd/gateway && sudo cp /tmp/mcp-gateway /usr/local/bin/mcp-gateway && sudo chmod +x /usr/local/bin/mcp-gateway && sudo mkdir -p /etc/mcp-gateway && sudo cp mapping.json /etc/mcp-gateway/ && sudo cp docker-catalog.yaml /etc/mcp-gateway/",
  ])

  // --- Clean up build artifacts ---
  .runCmd([
    "sudo rm -rf /opt/mcp-gateway-src /usr/local/go",
  ])

  .setWorkdir("/etc/mcp-gateway")

Template.build(template, 'claude-mcp', {
  cpuCount: 4,
  memoryMB: 8192,
  onBuildLogs: defaultBuildLogger(),
})

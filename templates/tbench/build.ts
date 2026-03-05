import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

const TEMPLATE_NAME = 'tbench'
const tag = process.env.E2B_BUILD_TAG || undefined

export const template = Template()
  .fromUbuntuImage("24.04")

  // --- System packages (tmux, git, asciinema) ---
  .runCmd([
    "sudo apt-get update && sudo DEBIAN_FRONTEND=noninteractive DEBCONF_NOWARNINGS=yes apt-get install -y tmux git asciinema"
  ])

  // --- Docker CE ---
  .runCmd([
    "curl -fsSL https://get.docker.com | sudo sh",
    "sudo systemctl daemon-reload",
    "sudo systemctl start docker",
    "sudo usermod -aG docker user",
  ])

  // --- uv + tools (harbor, terminal-bench) ---
  .runCmd([
    'curl -LsSf https://astral.sh/uv/install.sh | sh',
    'source $HOME/.local/bin/env && echo \'source "$HOME/.local/bin/env"\' >> ~/.bashrc && uv tool install harbor && uv tool install terminal-bench',
  ])

Template.build(template, tag ? `${TEMPLATE_NAME}:${tag}` : TEMPLATE_NAME, {
  cpuCount: 2,
  memoryMB: 4096,
  onBuildLogs: defaultBuildLogger(),
})

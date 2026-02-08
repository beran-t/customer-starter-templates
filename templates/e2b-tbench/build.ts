import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["tmux", "git"])
  .runCmd([
    "curl -fsSL https://get.docker.com | sudo sh",
    "sudo mkdir -p /etc/systemd/system/docker.service.d",
    "echo '[Service]\nExecStartPost=/bin/chmod 666 /var/run/docker.sock' | sudo tee /etc/systemd/system/docker.service.d/socket-perms.conf",
    "sudo systemctl daemon-reload",
    "sudo systemctl start docker",
    "curl -LsSf https://astral.sh/uv/install.sh | sh",
    "source $HOME/.local/bin/env",
    'echo \'source "$HOME/.local/bin/env"\' >> ~/.bashrc',
    "uv tool install harbor",
    "uv tool install terminal-bench",
    "harbor --version",
    "tb run --help"
  ])

Template.build(template, 'e2b-tbench', {
  cpuCount: 2,
  memoryMB: 4096,
  onBuildLogs: defaultBuildLogger(),
})

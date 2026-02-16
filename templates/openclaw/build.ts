import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "sed -i '/# If not running interactively/,/esac/d' ~/.bashrc",
    "curl -fsSL --proto '=https' --tlsv1.2 https://openclaw.ai/install-cli.sh | bash",
    "echo 'export PATH=\"$HOME/.openclaw/bin:$PATH\"' >> ~/.bashrc",
  ])

Template.build(template, 'openclaw', {
  cpuCount: 8,
  memoryMB: 8192,
  onBuildLogs: defaultBuildLogger(),
})

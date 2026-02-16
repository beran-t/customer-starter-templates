import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "curl -fsSL --proto '=https' --tlsv1.2 https://openclaw.ai/install-cli.sh | bash",
    "sudo find / -name openclaw -type f -perm /111 -not -path '*/node_modules/*' 2>/dev/null | head -1 | xargs -I{} sudo ln -sf {} /usr/local/bin/openclaw",
  ])

Template.build(template, 'openclaw', {
  cpuCount: 8,
  memoryMB: 8192,
  onBuildLogs: defaultBuildLogger(),
})

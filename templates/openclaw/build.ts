import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "curl -fsSL --proto '=https' --tlsv1.2 https://openclaw.ai/install-cli.sh | bash",
  ])

Template.build(template, 'openclaw', {
  cpuCount: 8,
  memoryMB: 8192,
  onBuildLogs: defaultBuildLogger(),
})

import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "curl -fsSL https://opencode.ai/install | bash",
    "source ~/.bashrc",
  ])

Template.build(template, 'opencode', {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

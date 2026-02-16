import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .runCmd([
    "npm i -g @openai/codex",
  ])

Template.build(template, 'codex', {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

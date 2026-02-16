import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .runCmd([
    "sed -i '/# If not running interactively/,/esac/d' ~/.bashrc",
    "source ~/.bashrc && npm i -g @openai/codex",
  ])

Template.build(template, 'codex', {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .runCmd([
    "npm i -g @openai/codex",
  ])

Template.build(template, 'codex', {
  onBuildLogs: defaultBuildLogger(),
})

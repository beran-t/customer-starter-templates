import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

const TEMPLATE_NAME = 'codex'
const tag = process.env.E2B_BUILD_TAG || undefined

export const template = Template()
  .fromBaseImage()
  .runCmd([
    "sudo npm i -g @openai/codex",
  ])

Template.build(template, tag ? `${TEMPLATE_NAME}:${tag}` : TEMPLATE_NAME, {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

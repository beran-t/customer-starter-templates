import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

const TEMPLATE_NAME = 'amp-code'
const tag = process.env.E2B_BUILD_TAG || undefined

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "sed -i '/# If not running interactively/,/esac/d' ~/.bashrc",
    "curl -fsSL https://ampcode.com/install.sh | bash",
    'echo \'export PATH="$HOME/.local/bin:$PATH"\' >> ~/.bashrc',
  ])

Template.build(template, tag ? `${TEMPLATE_NAME}:${tag}` : TEMPLATE_NAME, {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "sed -i '/# If not running interactively/,/esac/d' ~/.bashrc",
    "curl -fsSL https://ampcode.com/install.sh | bash",
    'echo \'export PATH="$HOME/.local/bin:$PATH"\' >> ~/.bashrc',
  ])

Template.build(template, 'amp-code', {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})

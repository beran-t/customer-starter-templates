import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

export const template = Template()
  .fromBaseImage()
  .aptInstall(["curl"])
  .runCmd([
    "curl -fsSL https://ampcode.com/install.sh | bash",
    'echo \'export PATH="$HOME/.local/bin:$PATH"\' >> ~/.bashrc',
  ])

Template.build(template, 'amp', {
  onBuildLogs: defaultBuildLogger(),
})

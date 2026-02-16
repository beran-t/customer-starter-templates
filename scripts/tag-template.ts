import "dotenv/config"
import { Template } from 'e2b'
import { readFileSync } from 'fs'
import { resolve } from 'path'

const templateName = process.argv[2]
if (!templateName) {
  console.error('Usage: npx tsx scripts/tag-template.ts <template-name>')
  process.exit(1)
}

const versionFile = resolve('templates', templateName, 'version')
let version: string
try {
  version = readFileSync(versionFile, 'utf-8').trim()
} catch {
  console.error(`Could not read version file: ${versionFile}`)
  process.exit(1)
}

if (!/^\d+\.\d+\.\d+$/.test(version)) {
  console.error(`Invalid version format "${version}" in ${versionFile}. Expected X.Y.Z`)
  process.exit(1)
}

const sourceTag = process.env.E2B_SOURCE_TAG || 'dev'
const tags = [`v${version}`, 'lts']

console.log(`Assigning tags ${JSON.stringify(tags)} to ${templateName}:${sourceTag}`)

await Template.assignTags(`${templateName}:${sourceTag}`, tags)

console.log(`Successfully tagged ${templateName}:${sourceTag} as ${tags.join(', ')}`)

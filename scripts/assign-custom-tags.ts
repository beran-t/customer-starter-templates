import "dotenv/config"
import { Template } from 'e2b'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'

const templateName = process.argv[2]
if (!templateName) {
  console.error('Usage: npx tsx scripts/assign-custom-tags.ts <template-name>')
  process.exit(1)
}

const tagsFile = resolve('templates', templateName, 'tags.json')

if (!existsSync(tagsFile)) {
  console.log(`No tags.json found for ${templateName}, skipping custom tags.`)
  process.exit(0)
}

interface CustomTag {
  name: string
  reference: string
  description: string
}

let tags: CustomTag[]
try {
  tags = JSON.parse(readFileSync(tagsFile, 'utf-8'))
} catch (err) {
  console.error(`Failed to parse ${tagsFile}: ${err}`)
  process.exit(1)
}

if (!Array.isArray(tags)) {
  console.error(`${tagsFile} must contain a JSON array.`)
  process.exit(1)
}

if (tags.length === 0) {
  console.log(`No custom tags defined for ${templateName}, skipping.`)
  process.exit(0)
}

const TAG_RE = /^[a-z0-9][a-z0-9._-]*$/
const RESERVED = ['lts', 'base']

for (const tag of tags) {
  if (!tag.name || !tag.reference || !tag.description) {
    console.error(`Invalid entry in ${tagsFile}: each object must have "name", "reference", and "description" fields.`)
    console.error(`Got: ${JSON.stringify(tag)}`)
    process.exit(1)
  }

  if (!TAG_RE.test(tag.name)) {
    console.error(`Invalid custom tag name "${tag.name}". Must match ${TAG_RE}`)
    process.exit(1)
  }

  if (RESERVED.includes(tag.name) || /^v\d/.test(tag.name) || /^dev(-|$)/.test(tag.name)) {
    console.error(`Custom tag "${tag.name}" conflicts with reserved tags (lts, base, v*, dev*).`)
    process.exit(1)
  }
}

for (const tag of tags) {
  const source = `${templateName}:${tag.reference}`
  console.log(`Assigning tag "${tag.name}" to ${source} (${tag.description})`)
  await Template.assignTags(source, [tag.name])
  console.log(`Successfully tagged ${source} as ${tag.name}`)
}

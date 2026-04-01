import "dotenv/config"
import { Template } from 'e2b'
import { readFileSync, writeFileSync, existsSync } from 'fs'
import { resolve } from 'path'

interface CustomTag {
  name: string
  reference: string
  description: string
}

const TAG_RE = /^[a-z0-9][a-z0-9._-]*$/
const RESERVED = ['lts', 'base']

function validateTagName(name: string) {
  if (!TAG_RE.test(name)) {
    console.error(`Invalid tag name "${name}". Must match ${TAG_RE}`)
    process.exit(1)
  }
  if (RESERVED.includes(name) || /^v\d+\.\d+\.\d+/.test(name) || /^dev(-|$)/.test(name)) {
    console.error(`Tag name "${name}" conflicts with reserved tags (lts, base, vX.Y.Z, dev*).`)
    process.exit(1)
  }
}

function tagsFilePath(templateName: string): string {
  const dir = resolve('templates', templateName)
  if (!existsSync(dir)) {
    console.error(`Template "${templateName}" not found at ${dir}`)
    process.exit(1)
  }
  return resolve(dir, 'tags.json')
}

function readTags(filePath: string): CustomTag[] {
  if (!existsSync(filePath)) return []
  try {
    const data = JSON.parse(readFileSync(filePath, 'utf-8'))
    if (!Array.isArray(data)) {
      console.error(`${filePath} must contain a JSON array.`)
      process.exit(1)
    }
    return data
  } catch (err) {
    console.error(`Failed to parse ${filePath}: ${err}`)
    process.exit(1)
  }
}

function writeTags(filePath: string, tags: CustomTag[]) {
  writeFileSync(filePath, JSON.stringify(tags, null, 2) + '\n')
}

const usage = `Usage: npx tsx scripts/manage-custom-tags.ts <command> <template> [options]

Commands:
  builds <template>                                          List successful builds and their tags (requires E2B_API_KEY)
  list   <template>                                          List custom tags from tags.json
  add    <template> --name <tag> --reference <ref> --description <desc>  Add a custom tag
  remove <template> --name <tag>                             Remove a custom tag`

const command = process.argv[2]
const templateName = process.argv[3]

if (!command || !templateName) {
  console.log(usage)
  process.exit(1)
}

function getArg(flag: string): string {
  const idx = process.argv.indexOf(flag)
  if (idx === -1 || idx + 1 >= process.argv.length) {
    console.error(`Missing required argument: ${flag}`)
    console.log(usage)
    process.exit(1)
  }
  return process.argv[idx + 1]
}

const filePath = tagsFilePath(templateName)

switch (command) {
  case 'builds': {
    const tags = await Template.getTags(templateName)

    // Group tags by buildId
    const tagsByBuild = new Map<string, string[]>()
    for (const t of tags) {
      const list = tagsByBuild.get(t.buildId) || []
      list.push(t.tag)
      tagsByBuild.set(t.buildId, list)
    }

    if (tagsByBuild.size === 0) {
      console.log(`No successful builds with tags found for ${templateName}.`)
    } else {
      console.log(`Successful builds for ${templateName}:\n`)
      for (const [buildId, buildTags] of tagsByBuild) {
        console.log(`  ${buildId}  tags: ${buildTags.join(', ')}`)
      }
    }
    break
  }

  case 'list': {
    const tags = readTags(filePath)
    if (tags.length === 0) {
      console.log(`No custom tags for ${templateName}.`)
    } else {
      console.log(`Custom tags for ${templateName}:\n`)
      for (const tag of tags) {
        console.log(`  ${tag.name} -> ${tag.reference}  (${tag.description})`)
      }
    }
    break
  }

  case 'add': {
    const name = getArg('--name')
    const reference = getArg('--reference')
    const description = getArg('--description')

    validateTagName(name)

    const tags = readTags(filePath)
    const existing = tags.find(t => t.name === name)
    if (existing) {
      console.log(`Updating existing tag "${name}": ${existing.reference} -> ${reference}`)
      existing.reference = reference
      existing.description = description
    } else {
      console.log(`Adding tag "${name}" -> ${reference}`)
      tags.push({ name, reference, description })
    }

    writeTags(filePath, tags)
    console.log(`Saved ${filePath}`)
    break
  }

  case 'remove': {
    const name = getArg('--name')
    const tags = readTags(filePath)
    const filtered = tags.filter(t => t.name !== name)
    if (filtered.length === tags.length) {
      console.error(`Tag "${name}" not found for ${templateName}.`)
      process.exit(1)
    }

    writeTags(filePath, filtered)
    console.log(`Removed tag "${name}" from ${templateName}.`)
    break
  }

  default:
    console.error(`Unknown command: ${command}`)
    console.log(usage)
    process.exit(1)
}

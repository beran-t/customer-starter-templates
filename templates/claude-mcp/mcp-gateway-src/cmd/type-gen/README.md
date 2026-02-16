# Type Generator

A tool for transforming Docker MCP catalog definitions into JSON schemas and instruction mappings.

## Overview

The type generator reads MCP catalog files (YAML) and produces:
- **`spec.json`**: Combined JSON schema for all services
- **`mapping.json`**: Instruction mappings for runtime configuration

## Usage

### Basic Usage
```bash
./type-gen
```
Uses default catalogs: `./docker-catalog.yaml`

### With Custom Options
```bash
./type-gen \
  -catalogs ./docker-catalog.yaml \
  -spec-output spec.json \
  -mapping-output mapping.json \
  -verbose
```

### CLI Flags
- `-catalogs`: Comma-separated list of catalog file paths (default: `./docker-catalog.yaml`)
- `-spec-output`: Output path for generated JSON schema (default: `spec.json`)
- `-mapping-output`: Output path for generated instruction mapping (default: `mapping.json`)
- `-verbose`: Enable verbose debug logging (default: `false`)

## Architecture

The codebase is an absolute mess

## Output Format

### spec.json
Combined JSON schema with beautified service names as top-level properties:
```json
{
  "type": "object",
  "properties": {
    "github": {
      "type": "object",
      "required": ["token"],
      "properties": {
        "token": { "type": "string" }
      }
    }
  }
}
```

### mapping.json
Instruction mappings for runtime configuration:
```json
{
  "github": {
    "server": "github-mcp-server",
    "type": ""
  },
  "github.token": {
    "server": "github-mcp-server",
    "type": "secret",
    "envName": "GITHUB_TOKEN"
  }
}
```


## Development

### Building
```bash
go build -o type-gen ./cmd/type-gen
```

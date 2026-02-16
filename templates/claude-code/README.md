# claude-code

A sandbox template for running Claude Code with MCP gateway support. Includes Docker, Node.js 22, Python 3, and a custom mcp-gateway binary compiled from source with a fix that enables MCP tool discovery.

## Template ID

`claude-code`

## What's Pre-installed

- **Docker CE** — container runtime with optimized daemon config
- **Node.js 22** — JavaScript runtime
- **Python 3** — with pip and venv
- **Poetry** — Python dependency management
- **uv / uvx** — fast Python package manager
- **Claude Code CLI** — Anthropic's coding assistant
- **mcp-gateway** — MCP gateway binary (compiled from source with `HasTools: true` fix)
- **jq** — JSON processor
- **git** — version control
- **build-essential** — C/C++ compiler toolchain

## Build Configuration

- **CPU**: 4 cores
- **Memory**: 8192 MB
- **Base image**: Ubuntu 25.04

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("claude-code", timeout=60)
try:
    result = sbx.commands.run("claude --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('claude-code', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('claude --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

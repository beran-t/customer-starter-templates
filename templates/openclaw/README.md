# openclaw

A sandbox template with the [OpenClaw](https://openclaw.ai) CLI pre-installed.

## Template ID

`openclaw`

## What's Pre-installed

- **curl** — HTTP client
- **OpenClaw CLI** — installed via the official install script

## Build Configuration

- **CPU**: 8 cores
- **Memory**: 8192 MB

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("openclaw", timeout=60)
try:
    result = sbx.commands.run("openclaw --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('openclaw', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('openclaw --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

## Building

This template uses the programmatic E2B Template builder:

```bash
npx tsx templates/openclaw/build.ts
```

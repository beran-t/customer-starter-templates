# opencode

A sandbox template with the [OpenCode](https://opencode.ai) CLI pre-installed.

## Template ID

`opencode`

## What's Pre-installed

- **OpenCode CLI** — installed via the official install script

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("opencode", timeout=60)
try:
    result = sbx.commands.run("opencode --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('opencode', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('opencode --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

## Building

This template uses the programmatic E2B Template builder:

```bash
npx tsx templates/opencode/build.ts
```

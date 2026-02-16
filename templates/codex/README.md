# codex

A sandbox template with the [OpenAI Codex CLI](https://github.com/openai/codex) pre-installed.

## Template ID

`codex`

## What's Pre-installed

- **Codex CLI** — OpenAI's coding agent (`@openai/codex`)

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("codex", timeout=60)
try:
    result = sbx.commands.run("codex --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('codex', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('codex --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

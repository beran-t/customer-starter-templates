# amp

A sandbox template with the [AMP](https://ampcode.com) CLI pre-installed.

## Template ID

`amp`

## What's Pre-installed

- **AMP CLI** — installed via the official install script

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("amp", timeout=60)
try:
    result = sbx.commands.run("amp --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('amp', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('amp --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

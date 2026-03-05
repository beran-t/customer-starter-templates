# tbench

A sandbox template with [Terminal Bench](https://github.com/terminal-bench/terminal-bench) and [Harbor](https://github.com/av/harbor) pre-installed.

## Template ID

`tbench`

## What's Pre-installed

- **Docker** — Container runtime with socket permissions configured for non-root use
- **tmux** — Terminal multiplexer
- **git** — Version control
- **asciinema** — Terminal session recorder
- **uv** — Fast Python package manager
- **Harbor** — AI tool orchestration (`harbor`)
- **Terminal Bench** — Terminal benchmarking tool (`tb`)

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox.create("tbench", timeout=60)
try:
    result = sbx.commands.run("harbor --version")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('tbench', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('harbor --version');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

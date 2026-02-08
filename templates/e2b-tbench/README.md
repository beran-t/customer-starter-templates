# e2b-tbench

A sandbox template for running [terminal-bench](https://github.com/terminal-bench/terminal-bench) benchmarks with Docker and Harbor support.

## Template ID

`e2b-tbench`

## What's Pre-installed

- **tmux** — terminal multiplexer
- **git** — version control
- **Docker** — container runtime (socket permissions configured for non-root access)
- **uv** — fast Python package manager
- **harbor** — installed via uv
- **terminal-bench (`tb`)** — installed via uv

## Build Configuration

- **CPU**: 2 cores
- **Memory**: 4096 MB

## Usage

### Python

```python
from e2b import Sandbox

sbx = Sandbox("e2b-tbench", timeout=60)
try:
    result = sbx.commands.run("tb run --help")
    print(result.stdout)
finally:
    sbx.kill()
```

### TypeScript

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('e2b-tbench', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('tb run --help');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

## Building

This template uses the programmatic E2B Template builder:

```bash
npx tsx templates/e2b-tbench/build.ts
```

# E2B Sandbox Templates

[![Test Templates](https://github.com/e2b-dev/e2b-sandbox-templates/actions/workflows/test-templates-pr.yml/badge.svg)](https://github.com/e2b-dev/e2b-sandbox-templates/actions/workflows/test-templates-pr.yml)

Pre-built sandbox templates for the [E2B](https://e2b.dev) platform. Each template defines a reproducible environment that users can spin up as a sandbox using the E2B SDK.

## Available Templates

| Name | Template ID | Description |
|------|-------------|-------------|
| [amp](templates/amp/) | `amp` | AMP CLI sandbox |
| [claude](templates/claude/) | `claude` | Claude Code sandbox with Docker, Node.js 22, Python 3, and MCP gateway |
| [codex](templates/codex/) | `codex` | OpenAI Codex CLI sandbox |
| [openclaw](templates/openclaw/) | `openclaw` | OpenClaw CLI sandbox |
| [opencode](templates/opencode/) | `opencode` | OpenCode CLI sandbox |
| [sandbox-egress-header](templates/sandbox-egress-header/) | `sandbox-egress-header` | Transparent proxy that injects X-Sandbox-ID header into all egress HTTP/HTTPS traffic |
| [tbench](templates/tbench/) | `tbench` | Terminal-bench sandbox with Docker, Harbor, and uv |

## Quickstart

Install the E2B SDK for your language:

```bash
# Python
pip install e2b

# TypeScript
npm install e2b
```

Use a template to create a sandbox:

**Python**

```python
from e2b import Sandbox

sbx = Sandbox.create("<template-name>", timeout=60)
try:
    result = sbx.commands.run("echo 'hello'")
    print(result.stdout)
finally:
    sbx.kill()
```

**TypeScript**

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('<template-name>', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('echo "hello"');
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for instructions on adding a new template.

## Documentation

- [E2B Documentation](https://e2b.dev/docs)
- [E2B Python SDK](https://pypi.org/project/e2b/)
- [E2B TypeScript SDK](https://www.npmjs.com/package/e2b)

## License

Apache 2.0 — see [LICENSE](LICENSE).

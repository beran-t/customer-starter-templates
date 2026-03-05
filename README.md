# Customer Starter Templates

![Test Templates](assets/test-templates-badge.svg)

A collection of pre-built [E2B](https://e2b.dev) sandbox templates shipped to customers as ready-to-use starting points. Each template defines a reproducible sandbox environment with the tools, runtimes, and configurations needed for a specific coding agent or workflow.

## Templates

| Template | Description |
|----------|-------------|
| [amp](templates/amp/) | AMP CLI sandbox |
| [claude](templates/claude/) | Claude Code sandbox with Docker, Node.js 22, Python 3, and MCP gateway |
| [codex](templates/codex/) | OpenAI Codex CLI sandbox |
| [openclaw](templates/openclaw/) | OpenClaw CLI sandbox |
| [opencode](templates/opencode/) | OpenCode CLI sandbox |
| [sandbox-egress-header](templates/sandbox-egress-header/) | Transparent proxy that injects X-Sandbox-ID header into all egress HTTP/HTTPS traffic |
| [tbench](templates/tbench/) | Terminal-bench sandbox with Docker, Harbor, and uv |

Each template has its own README with usage examples in Python and TypeScript.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for instructions on adding or updating templates.

## License

Apache 2.0 — see [LICENSE](LICENSE).

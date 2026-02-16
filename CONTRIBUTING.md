# Contributing

Thanks for your interest in contributing a new sandbox template!

## Adding a New Template

### 1. Create a template directory

Create a new directory under `templates/` with your template name:

```bash
mkdir templates/my-template
```

### 2. Add a build file

Create `templates/my-template/build.ts` (programmatic build using the E2B SDK) or `templates/my-template/e2b.Dockerfile` with the environment setup for your template.

### 3. Write example files

Create `templates/my-template/example.py`:

```python
from e2b import Sandbox

sbx = Sandbox.create("my-template", timeout=60)
try:
    result = sbx.commands.run("echo 'hello'")
    assert result.exit_code == 0
    # Add template-specific verification here
finally:
    sbx.kill()
```

Create `templates/my-template/example.ts`:

```typescript
import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('my-template', { timeoutMs: 60_000 });
try {
  const result = await sbx.commands.run('echo "hello"');
  if (result.exitCode !== 0) throw new Error('Command failed');
  // Add template-specific verification here
} finally {
  await sbx.kill();
}
```

### 4. Write a README

Create `templates/my-template/README.md` covering:

- What the template provides
- The template ID/name
- Usage examples in Python and TypeScript
- What's pre-installed in the sandbox

### 5. Test locally

Run both example files to verify they work:

```bash
python templates/my-template/example.py
npx tsx templates/my-template/example.ts
```

### 6. Open a PR

Open a pull request. CI will automatically build your template and run both example files.

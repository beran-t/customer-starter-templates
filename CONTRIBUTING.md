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

### 3. Add a test script

Create `templates/my-template/test.sh` — a bash script that verifies all custom-installed tools from the user's perspective:

```bash
#!/bin/bash
set -euo pipefail

echo "=== My Template Test ==="

echo "Checking my-tool..."
my-tool --version

echo ""
echo "All checks passed."
```

### 4. Write example files

The example files create a sandbox, upload `test.sh`, and run it. This keeps test logic in one place.

Create `templates/my-template/example.py`:

```python
import os
from e2b import Sandbox

test_dir = os.path.dirname(os.path.abspath(__file__))
with open(os.path.join(test_dir, "test.sh")) as f:
    test_script = f.read()

sbx = Sandbox.create("my-template", timeout=60)
try:
    sbx.files.write("/tmp/test.sh", test_script)
    result = sbx.commands.run("bash /tmp/test.sh")
    assert result.exit_code == 0, f"Test failed:\n{result.stderr}"
    print(result.stdout)
finally:
    sbx.kill()
```

Create `templates/my-template/example.ts`:

```typescript
import { Sandbox } from 'e2b';
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const testScript = readFileSync(join(__dirname, 'test.sh'), 'utf-8');

const sbx = await Sandbox.create('my-template', { timeoutMs: 60_000 });
try {
  await sbx.files.write('/tmp/test.sh', testScript);
  const result = await sbx.commands.run('bash /tmp/test.sh');
  if (result.exitCode !== 0) throw new Error(`Test failed:\n${result.stderr}`);
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
```

### 5. Write a README

Create `templates/my-template/README.md` covering:

- What the template provides
- The template ID/name
- Usage examples in Python and TypeScript
- What's pre-installed in the sandbox

### 6. Test locally

Run both example files to verify they work:

```bash
python templates/my-template/example.py
npx tsx templates/my-template/example.ts
```

### 7. Open a PR

Open a pull request. CI will automatically build your template and run both example files.

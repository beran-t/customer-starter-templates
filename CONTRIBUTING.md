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

The build file should support the `E2B_BUILD_TAG` env var for tagged builds:

```typescript
import "dotenv/config"
import { Template, defaultBuildLogger } from 'e2b'

const TEMPLATE_NAME = 'my-template'
const tag = process.env.E2B_BUILD_TAG || undefined

export const template = Template()
  .fromBaseImage()
  .runCmd(["my-tool install command"])

Template.build(template, tag ? `${TEMPLATE_NAME}:${tag}` : TEMPLATE_NAME, {
  cpuCount: 2,
  memoryMB: 2048,
  onBuildLogs: defaultBuildLogger(),
})
```

### 3. Add a version file

Create `templates/my-template/version` with the initial version:

```
1.0.0
```

Bump this version whenever you make changes to the template. CI will fail if you try to publish a version that already exists.

### 4. Write example files

The example files create a sandbox and verify all custom-installed tools work from the user's perspective. They should support the `E2B_TEMPLATE_TAG` env var for testing tagged builds.

Create `templates/my-template/example.py`:

```python
import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"my-template:{tag}" if tag else "my-template"

sbx = Sandbox.create(template_ref, timeout=60)
try:
    result = sbx.commands.run("my-tool --version")
    assert result.exit_code == 0, f"my-tool check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()
```

Create `templates/my-template/example.ts`:

```typescript
import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `my-template:${tag}` : 'my-template';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const myTool = await sbx.commands.run('my-tool --version');
  if (myTool.exitCode !== 0) throw new Error(`my-tool check failed: ${myTool.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
```

### 6. Write a README

Create `templates/my-template/README.md` covering:

- What the template provides
- The template ID/name
- Usage examples in Python and TypeScript
- What's pre-installed in the sandbox

### 7. Test locally

Run both example files to verify they work:

```bash
# Build and test without tags (backwards compatible)
npx tsx templates/my-template/build.ts
python templates/my-template/example.py
npx tsx templates/my-template/example.ts

# Or use the test runner with a tag
./scripts/run-tests.sh my-template dev
```

### 8. Open a PR

Open a pull request. CI will automatically build your template as `:dev` and run both example files against it.

## Versioning

Templates use a `dev` -> versioned -> `lts` promotion pipeline:

- **PR**: builds `template:dev`, runs tests against it
- **Merge to main**: builds `template:dev`, tests it, then tags with `vX.Y.Z` + `lts`

### Version file

Each template has a `version` file containing a semver string (e.g., `1.0.0`). Bump this when making changes.

### Tags

| Tag | Meaning |
|-----|---------|
| `dev` | Latest build, may not be tested yet |
| `vX.Y.Z` | Specific version, immutable |
| `lts` | Latest tested + promoted version |

Consumers should use `Sandbox.create('template:lts')` for stability.

### Environment variables

| Variable | Used by | Purpose |
|----------|---------|---------|
| `E2B_BUILD_TAG` | `build.ts` | Tag to assign when building (default: none) |
| `E2B_TEMPLATE_TAG` | `example.*` | Tag to use when creating sandbox (default: none) |
| `E2B_SOURCE_TAG` | `tag-template.ts` | Source tag to promote from (default: `dev`) |

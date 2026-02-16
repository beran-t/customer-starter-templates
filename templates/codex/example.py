import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"codex:{tag}" if tag else "codex"

sbx = Sandbox.create(template_ref, timeout=60)
try:
    result = sbx.commands.run("codex --version")
    assert result.exit_code == 0, f"codex check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

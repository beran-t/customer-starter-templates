import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"opencode:{tag}" if tag else "opencode"

sbx = Sandbox.create(template_ref, timeout=60)
try:
    result = sbx.commands.run("opencode --version")
    assert result.exit_code == 0, f"opencode check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

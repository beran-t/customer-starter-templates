import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"amp-code:{tag}" if tag else "amp-code"

sbx = Sandbox.create(template_ref, timeout=60)
try:
    result = sbx.commands.run("amp --version")
    assert result.exit_code == 0, f"amp check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

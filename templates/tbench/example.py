import os
from e2b import Sandbox

tag = os.environ.get("E2B_TEMPLATE_TAG", "")
template_ref = f"tbench:{tag}" if tag else "tbench"

sbx = Sandbox.create(template_ref, timeout=60)
try:
    result = sbx.commands.run("harbor --version")
    assert result.exit_code == 0, f"harbor check failed: {result.stderr}"

    result = sbx.commands.run("tb run --help")
    assert result.exit_code == 0, f"terminal-bench check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

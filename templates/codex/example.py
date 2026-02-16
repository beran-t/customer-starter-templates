from e2b import Sandbox

sbx = Sandbox.create("codex", timeout=60)
try:
    result = sbx.commands.run("codex --version")
    assert result.exit_code == 0, f"codex check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

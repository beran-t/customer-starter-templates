from e2b import Sandbox

sbx = Sandbox.create("openclaw", timeout=60)
try:
    result = sbx.commands.run("openclaw --version")
    assert result.exit_code == 0, f"openclaw check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

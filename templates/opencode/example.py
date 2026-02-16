from e2b import Sandbox

sbx = Sandbox.create("opencode", timeout=60)
try:
    result = sbx.commands.run("opencode --version")
    assert result.exit_code == 0, f"opencode check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

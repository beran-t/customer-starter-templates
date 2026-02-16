from e2b import Sandbox

sbx = Sandbox.create("amp", timeout=60)
try:
    result = sbx.commands.run("amp --version")
    assert result.exit_code == 0, f"amp check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

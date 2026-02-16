from e2b import Sandbox

sbx = Sandbox.create("claude-code", timeout=60)
try:
    result = sbx.commands.run("docker --version")
    assert result.exit_code == 0, f"docker check failed: {result.stderr}"

    result = sbx.commands.run("node --version")
    assert result.exit_code == 0, f"node check failed: {result.stderr}"

    result = sbx.commands.run("python3 --version")
    assert result.exit_code == 0, f"python3 check failed: {result.stderr}"

    result = sbx.commands.run("claude --version")
    assert result.exit_code == 0, f"claude check failed: {result.stderr}"

    result = sbx.commands.run("mcp-gateway --help")
    assert result.exit_code == 0, f"mcp-gateway check failed: {result.stderr}"

    result = sbx.commands.run("uv --version")
    assert result.exit_code == 0, f"uv check failed: {result.stderr}"

    result = sbx.commands.run("poetry --version")
    assert result.exit_code == 0, f"poetry check failed: {result.stderr}"

    result = sbx.commands.run("jq --version")
    assert result.exit_code == 0, f"jq check failed: {result.stderr}"

    result = sbx.commands.run("git --version")
    assert result.exit_code == 0, f"git check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

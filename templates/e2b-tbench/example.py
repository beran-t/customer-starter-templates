from e2b import Sandbox

sbx = Sandbox("e2b-tbench", timeout=60)
try:
    # Verify tmux is installed
    result = sbx.commands.run("tmux -V")
    assert result.exit_code == 0, f"tmux check failed: {result.stderr}"

    # Verify git is installed
    result = sbx.commands.run("git --version")
    assert result.exit_code == 0, f"git check failed: {result.stderr}"

    # Verify Docker is available
    result = sbx.commands.run("docker --version")
    assert result.exit_code == 0, f"docker check failed: {result.stderr}"

    # Verify uv is installed
    result = sbx.commands.run("uv --version")
    assert result.exit_code == 0, f"uv check failed: {result.stderr}"

    # Verify harbor is installed
    result = sbx.commands.run("harbor --version")
    assert result.exit_code == 0, f"harbor check failed: {result.stderr}"

    # Verify terminal-bench is installed
    result = sbx.commands.run("tb run --help")
    assert result.exit_code == 0, f"terminal-bench check failed: {result.stderr}"

    print("All checks passed.")
finally:
    sbx.kill()

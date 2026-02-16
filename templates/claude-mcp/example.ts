import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('claude-mcp', { timeoutMs: 60_000 });
try {
  // Verify Docker is available
  const docker = await sbx.commands.run('docker --version');
  if (docker.exitCode !== 0) throw new Error(`docker check failed: ${docker.stderr}`);

  // Verify Node.js is installed
  const node = await sbx.commands.run('node --version');
  if (node.exitCode !== 0) throw new Error(`node check failed: ${node.stderr}`);

  // Verify Python is installed
  const python = await sbx.commands.run('python3 --version');
  if (python.exitCode !== 0) throw new Error(`python3 check failed: ${python.stderr}`);

  // Verify Claude Code CLI is installed
  const claude = await sbx.commands.run('claude --version');
  if (claude.exitCode !== 0) throw new Error(`claude check failed: ${claude.stderr}`);

  // Verify mcp-gateway binary is installed
  const gateway = await sbx.commands.run('mcp-gateway --help');
  if (gateway.exitCode !== 0) throw new Error(`mcp-gateway check failed: ${gateway.stderr}`);

  // Verify uv is installed
  const uv = await sbx.commands.run('uv --version');
  if (uv.exitCode !== 0) throw new Error(`uv check failed: ${uv.stderr}`);

  // Verify poetry is installed
  const poetry = await sbx.commands.run('poetry --version');
  if (poetry.exitCode !== 0) throw new Error(`poetry check failed: ${poetry.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `claude-code:${tag}` : 'claude-code';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const docker = await sbx.commands.run('docker --version');
  if (docker.exitCode !== 0) throw new Error(`docker check failed: ${docker.stderr}`);

  const node = await sbx.commands.run('node --version');
  if (node.exitCode !== 0) throw new Error(`node check failed: ${node.stderr}`);

  const python = await sbx.commands.run('python3 --version');
  if (python.exitCode !== 0) throw new Error(`python3 check failed: ${python.stderr}`);

  const claude = await sbx.commands.run('claude --version');
  if (claude.exitCode !== 0) throw new Error(`claude check failed: ${claude.stderr}`);

  const gateway = await sbx.commands.run('mcp-gateway --help');
  if (gateway.exitCode !== 0) throw new Error(`mcp-gateway check failed: ${gateway.stderr}`);

  const uv = await sbx.commands.run('uv --version');
  if (uv.exitCode !== 0) throw new Error(`uv check failed: ${uv.stderr}`);

  const poetry = await sbx.commands.run('poetry --version');
  if (poetry.exitCode !== 0) throw new Error(`poetry check failed: ${poetry.stderr}`);

  const jq = await sbx.commands.run('jq --version');
  if (jq.exitCode !== 0) throw new Error(`jq check failed: ${jq.stderr}`);

  const git = await sbx.commands.run('git --version');
  if (git.exitCode !== 0) throw new Error(`git check failed: ${git.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

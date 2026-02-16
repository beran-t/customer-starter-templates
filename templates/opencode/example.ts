import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('opencode', { timeoutMs: 60_000 });
try {
  const opencode = await sbx.commands.run('opencode --version');
  if (opencode.exitCode !== 0) throw new Error(`opencode check failed: ${opencode.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

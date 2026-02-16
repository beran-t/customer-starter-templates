import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('codex', { timeoutMs: 60_000 });
try {
  const codex = await sbx.commands.run('codex --version');
  if (codex.exitCode !== 0) throw new Error(`codex check failed: ${codex.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

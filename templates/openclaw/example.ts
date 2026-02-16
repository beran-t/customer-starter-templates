import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('openclaw', { timeoutMs: 60_000 });
try {
  const openclaw = await sbx.commands.run('openclaw --version');
  if (openclaw.exitCode !== 0) throw new Error(`openclaw check failed: ${openclaw.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}

import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('amp-code', { timeoutMs: 60_000 });
try {
  const amp = await sbx.commands.run('amp --version');
  if (amp.exitCode !== 0) throw new Error(`amp check failed: ${amp.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
